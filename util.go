package sqlutil

import "database/sql"

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

func execSQL(db *sql.DB, statement string, values ...interface{}) (int64, error) {
	result, err := db.Exec(statement, values...)
	if err != nil {
		return 0, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
