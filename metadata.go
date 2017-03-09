package sqlutil

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var metadata *Metadata

func init() {
	metadata = &Metadata{}
}

const (
	TagName               = "sql"
	TagSuffixIndexName    = "index"
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
	ColumnConstraintUnique ColumnConstraint = 1 << iota
	ColumnConstraintNull
	ColumnConstraintNotNull
)

type Schema struct {
	PrimaryKey []string
	Columns    []*Column
	Indexes    []*Index
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
		PrimaryKey: []string{},
		Columns:    []*Column{},
		Indexes:    []*Index{},
	}

	m.info[t] = schema
	indexes := map[string]*Index{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(TagName)

		if field.PkgPath != "" || tag == "-" {
			continue
		}

		if tag == "" {
			return nil, fmt.Errorf("Missing tag for field %q in type %q", field.Name, t.Name())
		}

		column := &Column{
			Index: i,
		}
		hasPrimaryKey := false

		for index, meta := range strings.Split(tag, ",") {
			switch index {
			case TagFieldNameIndex:
				column.Name = meta
			case TagFieldDataTypeIndex:
				column.DataType = meta
			default:
				if meta == "pk" {
					hasPrimaryKey = true
				} else if strings.HasPrefix(meta, TagSuffixIndexName) {
					m.index(column.Name, meta, indexes)
				} else {
					column.Constraint |= m.constraints(meta)
				}
			}
		}

		schema.Columns = append(schema.Columns, column)

		if hasPrimaryKey {
			schema.PrimaryKey = append(schema.PrimaryKey, column.Name)
		}
	}

	for _, index := range indexes {
		schema.Indexes = append(schema.Indexes, index)
	}

	sort.Slice(schema.Indexes, func(i, j int) bool {
		return schema.Indexes[i].Name < schema.Indexes[j].Name
	})

	return schema, nil
}

func (m *Metadata) index(columnName, meta string, indexes map[string]*Index) {
	indexName := fmt.Sprintf("%s_idx", columnName)
	parts := strings.Split(meta, ":")

	if len(parts) > 1 {
		indexName = parts[1]
	}

	if index, ok := indexes[indexName]; ok {
		index.Columns = append(index.Columns, columnName)
		return
	}

	indexes[indexName] = &Index{
		Name:    indexName,
		Columns: []string{columnName},
	}
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

func typeOf(m interface{}) (reflect.Type, error) {
	v := reflect.ValueOf(m)
	t := v.Type()

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("Must be pointer to struct; got %T", v)
	}

	return t.Elem(), nil
}

func valueOf(m interface{}) reflect.Value {
	v := reflect.ValueOf(m)
	return v.Elem()
}
