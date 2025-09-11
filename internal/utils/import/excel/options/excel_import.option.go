// ...existing code...
package options

import "strings"

type ExcelImportOption struct {
	HeaderRow      int
	AllSheet       bool
	SheetName      *string
	ValidationPath string
}

type Option func(*ExcelImportOption)

func WithHeaderRow(n int) Option {
	return func(o *ExcelImportOption) {
		if n > 0 {
			o.HeaderRow = n
		}
	}
}

func WithAllSheet(all bool) Option {
	return func(o *ExcelImportOption) {
		o.AllSheet = all
	}
}

func WithSheetName(name string) Option {
	return func(o *ExcelImportOption) {
		n := strings.TrimSpace(name)
		if n == "" {
			o.SheetName = nil
			return
		}
		o.SheetName = &n
	}
}

func WithValidationPath(path string) Option {
	return func(o *ExcelImportOption) {
		o.ValidationPath = path
	}
}

func DefaultOptions() *ExcelImportOption {
	return &ExcelImportOption{
		HeaderRow:      1,
		AllSheet:       false,
		SheetName:      nil,
		ValidationPath: "",
	}
}
