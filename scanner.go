package sqlutil

import "database/sql"

type RowScanner struct {
	Row *sql.Row
}

func (s *RowScanner) Scan(dest ...interface{}) error {
	return s.Row.Scan(dest...)
}

func (s *RowScanner) Columns() ([]string, error) {
	return []string{}, nil
}

type Scanner interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
}

func Scan(scanner Scanner, model interface{}) error {
	return NewEntityContext(model).Scan(scanner)
}
