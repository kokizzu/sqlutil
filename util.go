package sqlutil

import "database/sql"

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
