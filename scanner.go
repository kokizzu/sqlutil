package sqlutil

import "database/sql"

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

	mapping := make(map[string]int)
	values := make([]interface{}, 0)

	for _, c := range schema.Columns {
		mapping[c.Name] = c.Index
	}

	columns, err := scanner.Columns()
	if err != nil {
		return err
	}

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
