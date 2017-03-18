package sqlutil

import "strings"

const (
	ColumnConstraintUnique ColumnConstraint = 1 << iota
	ColumnConstraintNull
	ColumnConstraintNotNull
)

type Schema struct {
	Table       string
	ForeignKeys []*ForeignKey
	Columns     []*Column
	Indexes     []*Index
}

type ForeignKey struct {
	Columns               []string
	ReferenceTable        string
	ReferenceTableColumns []string
}

type Column struct {
	Name       string
	Index      int
	DataType   string
	PrimaryKey bool
	Constraint ColumnConstraint
}

type Index struct {
	Name    string
	Columns []string
}

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
