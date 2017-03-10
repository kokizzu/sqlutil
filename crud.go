package sqlutil

import (
	"database/sql"
	"fmt"
	"strings"
)

func Insert(db *sql.DB, model interface{}) error {
	t, err := typeOf(model)
	if err != nil {
		return err
	}

	schema, err := metadata.Schema(t)
	if err != nil {
		return err
	}

	v := valueOf(model)
	columns := []string{}
	values := make([]interface{}, 0)
	placeholders := []string{}

	for _, column := range schema.Columns {
		value := v.Field(column.Index).Addr().Interface()
		values = append(values, value)
		columns = append(columns, column.Name)
		placeholders = append(placeholders, "?")
	}

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", schema.Table, strings.Join(columns, ","), strings.Join(placeholders, ","))
	_, err = db.Exec(statement, values...)
	return err
}
