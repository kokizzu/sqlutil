package sqlutil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type TableEntity struct {
	schema     *Schema
	modelValue reflect.Value
}

func Entity(model interface{}) *TableEntity {
	typ, err := typeOf(model)
	if err != nil {
		panic(err)
	}

	schema, err := metadata.Schema(typ)
	if err != nil {
		panic(err)
	}

	return &TableEntity{
		modelValue: valueOf(model),
		schema:     schema,
	}
}

func (t *TableEntity) Scan(scanner Scanner) error {
	columns, err := scanner.Columns()
	if err != nil {
		return err
	}

	add := len(columns) == 0

	mapping := make(map[string]int)
	for _, c := range t.schema.Columns {
		mapping[c.Name] = c.Index
		if add {
			columns = append(columns, c.Name)
		}
	}

	values := make([]interface{}, 0)
	var value interface{}

	for _, column := range columns {
		if idx, ok := mapping[column]; ok {
			value = t.modelValue.Field(idx).Addr().Interface()
		} else {
			value = &sql.RawBytes{}
		}

		values = append(values, value)
	}

	return scanner.Scan(values...)
}

func (t *TableEntity) QueryRow(db *sql.DB) error {
	columns := []string{}
	values := make([]interface{}, 0)

	for _, column := range t.schema.Columns {
		value := t.modelValue.Field(column.Index).Addr().Interface()
		expression := fmt.Sprintf("%s = ?", column.Name)

		if column.Constraint&ColumnConstraintPrimaryKey != 0 {
			columns = append(columns, expression)
			values = append(values, value)
		}
	}

	statement := fmt.Sprintf("SELECT * FROM %s WHERE %s", t.schema.Table, strings.Join(columns, ","))
	row := db.QueryRow(statement, values...)
	return t.Scan(&RowScanner{row})
}

func (t *TableEntity) Insert(db *sql.DB) (int64, error) {
	columns := []string{}
	values := make([]interface{}, 0)
	placeholders := []string{}

	for _, column := range t.schema.Columns {
		value := t.modelValue.Field(column.Index).Addr().Interface()
		values = append(values, value)
		columns = append(columns, column.Name)
		placeholders = append(placeholders, "?")
	}

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", t.schema.Table, strings.Join(columns, ","), strings.Join(placeholders, ","))
	return execSQL(db, statement, values...)
}

func (t *TableEntity) Update(db *sql.DB) (int64, error) {
	columns := []string{}
	values := make([]interface{}, 0)
	conditionValues := make([]interface{}, 0)
	conditions := []string{}

	for _, column := range t.schema.Columns {
		value := t.modelValue.Field(column.Index).Addr().Interface()
		expression := fmt.Sprintf("%s = ?", column.Name)

		if column.Constraint&ColumnConstraintPrimaryKey != 0 {
			conditions = append(conditions, expression)
			conditionValues = append(values, value)
		} else {
			columns = append(columns, expression)
			values = append(values, value)
		}
	}

	values = append(values, conditionValues...)
	statement := fmt.Sprintf("UPDATE %s SET %s WHERE %s", t.schema.Table, strings.Join(columns, ","), strings.Join(conditions, ","))
	return execSQL(db, statement, values...)
}

func (t *TableEntity) Delete(db *sql.DB) (int64, error) {
	columns := []string{}
	values := make([]interface{}, 0)

	for _, column := range t.schema.Columns {
		value := t.modelValue.Field(column.Index).Addr().Interface()
		expression := fmt.Sprintf("%s = ?", column.Name)

		if column.Constraint&ColumnConstraintPrimaryKey != 0 {
			columns = append(columns, expression)
			values = append(values, value)
		}
	}

	statement := fmt.Sprintf("DELETE FROM %s WHERE %s", t.schema.Table, strings.Join(columns, ","))
	return execSQL(db, statement, values...)
}
