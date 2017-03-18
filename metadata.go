package sqlutil

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	metadata         *Metadata
	ignoredFieldErr  error = fmt.Errorf("Field is ignored")
	foreignKeyRegexp       = regexp.MustCompile(`([\w]+)\(([\w]+)\)`)
)

func init() {
	metadata = &Metadata{}
}

const (
	TagColumnName         = "sql"
	TagIndexName          = "sqlindex"
	TagForeignKeyName     = "sqlforeignkey"
	TagFieldNameIndex     = 0
	TagFieldDataTypeIndex = 1
)

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
		Table:       strings.ToLower(t.Name()),
		ForeignKeys: []*ForeignKey{},
		Columns:     []*Column{},
		Indexes:     []*Index{},
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
		m.foreignKey(schema, column, field)

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
		if meta == "pk" {
			column.PrimaryKey = true
		} else {
			switch index {
			case TagFieldNameIndex:
				column.Name = meta
			case TagFieldDataTypeIndex:
				column.DataType = meta
			default:
				column.Constraint |= m.constraints(meta)
			}
		}
	}

	return nil
}

func (m *Metadata) constraints(meta string) ColumnConstraint {
	switch meta {
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

func (m *Metadata) foreignKey(schema *Schema, column *Column, field reflect.StructField) {
	tag := Tag(field.Tag)

	for _, fkTag := range tag.Get(TagForeignKeyName) {
		found := false

		matches := foreignKeyRegexp.FindStringSubmatch(fkTag)
		if len(matches) < 3 {
			continue
		}

		referenceTable := matches[1]
		referenceColumn := matches[2]

		for _, fk := range schema.ForeignKeys {
			if fk.ReferenceTable == referenceTable {
				fk.ReferenceTableColumns = append(fk.ReferenceTableColumns, referenceColumn)
				fk.Columns = append(fk.Columns, column.Name)
				found = true
				break
			}
		}

		if !found {
			schema.ForeignKeys = append(schema.ForeignKeys, &ForeignKey{
				ReferenceTable:        referenceTable,
				ReferenceTableColumns: []string{referenceColumn},
				Columns:               []string{column.Name},
			})
		}
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
