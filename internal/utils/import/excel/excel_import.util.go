package utilsImportExcel

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"movie-app-go/internal/utils/import/excel/options"

	"github.com/xuri/excelize/v2"
	"slices"
)

func ParseSheetToMapsWithHeader(r io.Reader, opts ...options.Option) ([]map[string]string, error) {
	opt := options.DefaultOptions()
	for _, o := range opts {
		o(opt)
	}

	if opt.HeaderRow <= 0 {
		return nil, fmt.Errorf("headerRow harus > 0")
	}
	if strings.TrimSpace(opt.ValidationPath) == "" {
		return nil, fmt.Errorf("validationPath wajib di-set pada options")
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file upload: %w", err)
	}

	if err := ValidateFile(bytes.NewReader(buf), opt.ValidationPath, opt.HeaderRow, opt.AllSheet, opt.SheetName); err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("tidak ada sheet pada file excel")
	}

	var targetSheets []string
	if opt.AllSheet {
		targetSheets = sheets
	} else if opt.SheetName != nil {
		found := slices.Contains(sheets, *opt.SheetName)
		if !found {
			return nil, fmt.Errorf("sheet '%s' tidak ditemukan dalam file", *opt.SheetName)
		}
		targetSheets = []string{*opt.SheetName}
	} else {
		targetSheets = []string{sheets[0]}
	}

	var results []map[string]string
	for _, sheetName := range targetSheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("gagal baca row sheet '%s': %w", sheetName, err)
		}
		if len(rows) < opt.HeaderRow {
			continue
		}
		headers := rows[opt.HeaderRow-1]
		for i, row := range rows {
			if i < opt.HeaderRow {
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
	}

	return results, nil
}

func ParseMultipartFileWithOptions(fh *multipart.FileHeader, opts ...options.Option) ([]map[string]string, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file upload: %w", err)
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file upload: %w", err)
	}

	return ParseSheetToMapsWithHeader(bytes.NewReader(buf), opts...)
}

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

func ValidateFile(uploadedFile io.Reader, templatePath string, headerRow int, allSheet bool, sheetName *string) error {
    tplF, err := os.Open(templatePath)
    if err != nil {
        return fmt.Errorf("gagal membuka template file: %w", err)
    }
    defer tplF.Close()

    tf, err := excelize.OpenReader(tplF)
    if err != nil {
        return fmt.Errorf("gagal membuka template excel: %w", err)
    }
    defer tf.Close()
    templateSheets := tf.GetSheetList()
    if len(templateSheets) == 0 {
        return fmt.Errorf("template tidak memiliki sheet")
    }

	templateFirst := templateSheets[0]
    trows, err := tf.GetRows(templateFirst)
    if err != nil {
        return fmt.Errorf("gagal baca sheet template '%s': %w", templateFirst, err)
    }
    if len(trows) < headerRow {
        return fmt.Errorf("template tidak memiliki header di row %d", headerRow)
    }
    templateHeaders := trows[headerRow-1]

    uf, err := excelize.OpenReader(uploadedFile)
    if err != nil {
        return fmt.Errorf("gagal membuka file upload: %w", err)
    }
    defer uf.Close()
    uploadSheets := uf.GetSheetList()
    if len(uploadSheets) == 0 {
        return fmt.Errorf("file upload tidak memiliki sheet")
    }

    compareHeaders := func(uSheet string) error {
        urows, err := uf.GetRows(uSheet)
        if err != nil {
            return fmt.Errorf("gagal baca sheet upload '%s': %w", uSheet, err)
        }
        if len(urows) < headerRow {
            return fmt.Errorf("file upload tidak memiliki header di sheet '%s' row %d", uSheet, headerRow)
        }
        uploadHeaders := urows[headerRow-1]

        if len(uploadHeaders) != len(templateHeaders) {
            return fmt.Errorf("jumlah kolom sheet '%s' tidak sesuai template. Expected: %d, Got: %d", uSheet, len(templateHeaders), len(uploadHeaders))
        }
        for i, expected := range templateHeaders {
            actual := strings.TrimSpace(uploadHeaders[i])
            if !strings.EqualFold(strings.TrimSpace(expected), actual) {
                return fmt.Errorf("header kolom %d sheet '%s' tidak sesuai template. Expected: '%s', Got: '%s'", i+1, uSheet, expected, actual)
            }
        }
        return nil
    }

    if sheetName != nil {
        foundUp := slices.Contains(uploadSheets, *sheetName)
        if !foundUp {
            return fmt.Errorf("file upload tidak memiliki sheet '%s'", *sheetName)
        }
        return compareHeaders(*sheetName)
    }

    if allSheet {
        for _, us := range uploadSheets {
            if err := compareHeaders(us); err != nil {
                return err
            }
        }
        return nil
    }

    return compareHeaders(uploadSheets[0])
}
