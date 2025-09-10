package utilsImportExcel

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetHeaders(r io.Reader, sheetName string, headerRow int) ([]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("tidak ada sheet pada file excel")
	}
	if sheetName == "" {
		sheetName = sheets[0]
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal baca row: %w", err)
	}
	if len(rows) < headerRow {
		return nil, fmt.Errorf("file tidak memiliki baris header yang diminta")
	}

	return rows[headerRow-1], nil
}

func ParseSheetToMapsWithHeader(r io.Reader, sheetName string, headerRow int) ([]map[string]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("tidak ada sheet pada file excel")
	}
	if sheetName == "" {
		sheetName = sheets[0]
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal baca row: %w", err)
	}
	if len(rows) < headerRow {
		return nil, fmt.Errorf("file tidak memiliki baris header yang diminta")
	}

	headers := rows[headerRow-1]
	var results []map[string]string
	for i, row := range rows {
		if i < headerRow {
			continue
		}

		empty := true
		for _, c := range row {
			if c != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}

		m := map[string]string{}
		for idx, h := range headers {
			val := ""
			if idx < len(row) {
				val = row[idx]
			}
			m[h] = val
		}
		results = append(results, m)
	}
	return results, nil
}

func ValidateFileHeadersWithTemplate(uploadedFile io.Reader, templatePath string, headerRow int) error {
	templateFile, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("gagal membuka template file: %w", err)
	}
	defer templateFile.Close()

	templateHeaders, err := GetHeaders(templateFile, "", headerRow)
	if err != nil {
		return fmt.Errorf("gagal membaca header template: %w", err)
	}

	uploadHeaders, err := GetHeaders(uploadedFile, "", headerRow)
	if err != nil {
		return fmt.Errorf("gagal membaca header file upload: %w", err)
	}

	if len(uploadHeaders) != len(templateHeaders) {
		return fmt.Errorf("jumlah kolom tidak sesuai template. Expected: %d, Got: %d", len(templateHeaders), len(uploadHeaders))
	}

	for i, expectedHeader := range templateHeaders {
		actualHeader := strings.TrimSpace(uploadHeaders[i])
		expectedHeader = strings.TrimSpace(expectedHeader)

		if !strings.EqualFold(actualHeader, expectedHeader) {
			return fmt.Errorf("header kolom %d tidak sesuai template. Expected: '%s', Got: '%s'", i+1, expectedHeader, actualHeader)
		}
	}

	return nil
}

func ValidateFileHeadersWithTemplateFromPath(uploadedFilePath, templatePath string, headerRow int) error {
	uploadedFile, err := os.Open(uploadedFilePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file upload: %w", err)
	}
	defer uploadedFile.Close()

	return ValidateFileHeadersWithTemplate(uploadedFile, templatePath, headerRow)
}

func ValidateMultipartFileWithTemplate(fileHeader *multipart.FileHeader, templatePath string, headerRow int) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file upload: %w", err)
	}
	defer file.Close()

	return ValidateFileHeadersWithTemplate(file, templatePath, headerRow)
}
