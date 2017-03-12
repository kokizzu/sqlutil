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
	t, err := typeOf(model)
	if err != nil {
		return err
	}

	v := valueOf(model)

	schema, err := metadata.Schema(t)
	if err != nil {
		return err
	}

	columns, err := scanner.Columns()
	if err != nil {
		return err
	}

	add := len(columns) == 0

	mapping := make(map[string]int)
	for _, c := range schema.Columns {
		mapping[c.Name] = c.Index
		if add {
			columns = append(columns, c.Name)
		}
	}

	values := make([]interface{}, 0)
	var value interface{}

	for _, column := range columns {
		if idx, ok := mapping[column]; ok {
			value = v.Field(idx).Addr().Interface()
		} else {
			value = &sql.RawBytes{}
		}

		values = append(values, value)
	}

	return scanner.Scan(values...)
}
