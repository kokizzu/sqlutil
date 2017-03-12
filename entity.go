package sqlutil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
)

type Fields map[string]interface{}

type EntityContext struct {
	schema     *Schema
	modelValue reflect.Value
}

func Entity(model interface{}) *EntityContext {
	typ, err := typeOf(model)
	if err != nil {
		panic(err)
	}

	schema, err := metadata.Schema(typ)
	if err != nil {
		panic(err)
	}

	return &EntityContext{
		modelValue: valueOf(model),
		schema:     schema,
	}
}

func (t *EntityContext) Scan(scanner Scanner) error {
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

func (t *EntityContext) QueryRow(db *sql.DB) error {
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

func (t *EntityContext) Insert(db *sql.DB) (int64, error) {
	columns := []string{}
	values := make([]interface{}, 0)
	placeholders := []string{}
	now := reflect.ValueOf(time.Now())

	for _, column := range t.schema.Columns {
		field := t.modelValue.Field(column.Index)
		if column.Name == FieldCreatedAt || column.Name == FieldUpdatedAt {
			field.Set(now)
		}

		value := field.Addr().Interface()
		values = append(values, value)
		columns = append(columns, column.Name)
		placeholders = append(placeholders, "?")
	}

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", t.schema.Table, strings.Join(columns, ","), strings.Join(placeholders, ","))
	return execSQL(db, statement, values...)
}

func (t *EntityContext) Update(db *sql.DB, fields ...Fields) (int64, error) {
	columns := []string{}
	values := make([]interface{}, 0)
	conditionValues := make([]interface{}, 0)
	conditions := []string{}
	allFields, merged := mergeFields(fields)
	now := reflect.ValueOf(time.Now())

	for _, column := range t.schema.Columns {
		field := t.modelValue.Field(column.Index)
		if column.Name == FieldUpdatedAt {
			field.Set(now)
		}

		value := field.Addr().Interface()
		expression := fmt.Sprintf("%s = ?", column.Name)

		if column.Constraint&ColumnConstraintPrimaryKey != 0 {
			conditions = append(conditions, expression)
			conditionValues = append(values, value)
			continue
		}

		if merged {
			ok := false
			if value, ok = allFields[column.Name]; !ok {
				continue
			}
		}

		columns = append(columns, expression)
		values = append(values, value)
	}

	values = append(values, conditionValues...)
	statement := fmt.Sprintf("UPDATE %s SET %s WHERE %s", t.schema.Table, strings.Join(columns, ","), strings.Join(conditions, ","))
	return execSQL(db, statement, values...)
}

func (t *EntityContext) Delete(db *sql.DB) (int64, error) {
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
