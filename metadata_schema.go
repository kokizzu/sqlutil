package sqlutil

import "strings"

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
