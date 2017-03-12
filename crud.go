package sqlutil

import (
	"database/sql"
)

func QueryRow(db *sql.DB, model interface{}) error {
	return Model(model).QueryRow(db)
}

func Insert(db *sql.DB, model interface{}) (int64, error) {
	return Model(model).Insert(db)
}

func Update(db *sql.DB, model interface{}) (int64, error) {
	return Model(model).Update(db)
}

func Delete(db *sql.DB, model interface{}) (int64, error) {
	return Model(model).Delete(db)
}
