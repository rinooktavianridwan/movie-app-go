package exportExcelUtil

import (
    "bytes"
    "fmt"

    "github.com/xuri/excelize/v2"
)

func GenerateFromMaps(sheetName string, headers []string, records []map[string]string) ([]byte, error) {
    f := excelize.NewFile()
    if sheetName == "" {
        sheetName = "Sheet1"
    }
    f.SetSheetName("Sheet1", sheetName)

    if err := f.SetSheetRow(sheetName, "A1", &headers); err != nil {
        return nil, fmt.Errorf("gagal set header: %w", err)
    }

    for i, rec := range records {
        row := make([]interface{}, len(headers))
        for j, h := range headers {
            row[j] = rec[h]
        }
        cell, _ := excelize.CoordinatesToCellName(1, i+2)
        if err := f.SetSheetRow(sheetName, cell, &row); err != nil {
            return nil, fmt.Errorf("gagal set row %d: %w", i+2, err)
        }
    }

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, fmt.Errorf("gagal write excel: %w", err)
    }
    return buf.Bytes(), nil
}