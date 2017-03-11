package sqlutil

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	metadata        *Metadata
	ignoredFieldErr error = fmt.Errorf("Field is ignored")
)

func init() {
	metadata = &Metadata{}
}

const (
	TagColumnName         = "sql"
	TagIndexName          = "sqlindex"
	TagFieldNameIndex     = 0
	TagFieldDataTypeIndex = 1
)

type ColumnConstraint byte

func (c ColumnConstraint) String() string {
	constraints := []string{}

	if c&ColumnConstraintUnique != 0 {
		constraints = append(constraints, "UNIQUE")
	}

	if c&ColumnConstraintNull != 0 {
		constraints = append(constraints, "NULL")
	}

	if c&ColumnConstraintNotNull != 0 {
		constraints = append(constraints, "NOT NULL")
	}

	return strings.Join(constraints, " ")
}

const (
	ColumnConstraintPrimaryKey ColumnConstraint = 1 << iota
	ColumnConstraintUnique
	ColumnConstraintNull
	ColumnConstraintNotNull
)

type Schema struct {
	Table   string
	Columns []*Column
	Indexes []*Index
}

type Column struct {
	Name       string
	Index      int
	DataType   string
	Constraint ColumnConstraint
}

type Index struct {
	Name    string
	Columns []string
}

type Metadata struct {
	info map[reflect.Type]*Schema
}

func (m *Metadata) Schema(t reflect.Type) (*Schema, error) {
	if m.info == nil {
		m.info = map[reflect.Type]*Schema{}
	}

	schema, ok := m.info[t]
	if ok {
		return schema, nil
	}

	schema = &Schema{
		Table:   strings.ToLower(t.Name()),
		Columns: []*Column{},
		Indexes: []*Index{},
	}

	m.info[t] = schema

	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)

		if field.PkgPath != "" {
			continue
		}

		column := &Column{
			Index: index,
		}

		if err := m.column(column, field); err != nil {
			if err == ignoredFieldErr {
				continue
			}
			return nil, fmt.Errorf("Type %q: %v", t.Name(), err)
		}

		m.index(schema, column, field)

		schema.Columns = append(schema.Columns, column)
	}

	return schema, nil
}

func (m *Metadata) column(column *Column, field reflect.StructField) error {
	columnTag := field.Tag.Get(TagColumnName)

	if columnTag == "-" {
		return ignoredFieldErr
	}

	if columnTag == "" {
		return fmt.Errorf("Missing tag for field %q", field.Name)
	}

	for index, meta := range strings.Split(columnTag, ",") {
		switch index {
		case TagFieldNameIndex:
			column.Name = meta
		case TagFieldDataTypeIndex:
			column.DataType = meta
		default:
			column.Constraint |= m.constraints(meta)
		}
	}

	return nil
}

func (m *Metadata) constraints(meta string) ColumnConstraint {
	switch meta {
	case "pk":
		return ColumnConstraintPrimaryKey
	case "unique":
		return ColumnConstraintUnique
	case "not_null":
		return ColumnConstraintNotNull
	case "null":
		return ColumnConstraintNull
	default:
		return ColumnConstraint(0)
	}
}

func (m *Metadata) index(schema *Schema, column *Column, field reflect.StructField) {
	tag := Tag(field.Tag)

	for _, indexTag := range tag.Get(TagIndexName) {
		found := false

		for _, index := range schema.Indexes {
			if index.Name == indexTag {
				index.Columns = append(index.Columns, column.Name)
				found = true
				break
			}
		}

		if !found {
			schema.Indexes = append(schema.Indexes, &Index{
				Name:    indexTag,
				Columns: []string{column.Name},
			})
		}
	}
}

func typeOf(m interface{}) (reflect.Type, error) {
	v := reflect.ValueOf(m)
	t := v.Type()

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("Must be pointer to struct; got %s", t.Name())
	}

	return t.Elem(), nil
}

func valueOf(m interface{}) reflect.Value {
	v := reflect.ValueOf(m)
	return v.Elem()
}
