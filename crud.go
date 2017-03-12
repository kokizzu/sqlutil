package sqlutil

import (
	"database/sql"
)

func QueryRow(db *sql.DB, model interface{}) error {
	return Entity(model).QueryRow(db)
}

func Insert(db *sql.DB, model interface{}) (int64, error) {
	return Entity(model).Insert(db)
}

func Update(db *sql.DB, model interface{}) (int64, error) {
	return Entity(model).Update(db)
}

func Delete(db *sql.DB, model interface{}) (int64, error) {
	return Entity(model).Delete(db)
}
