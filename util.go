package sqlutil

import "database/sql"

func QueryRow(db *sql.DB, model interface{}) error {
	return NewEntityContext(model).QueryRow(db)
}

func Insert(db *sql.DB, model interface{}) (int64, error) {
	return NewEntityContext(model).Insert(db)
}

func Update(db *sql.DB, model interface{}) (int64, error) {
	return NewEntityContext(model).Update(db)
}

func Delete(db *sql.DB, model interface{}) (int64, error) {
	return NewEntityContext(model).Delete(db)
}

func mergeFields(fields []Fields) (Fields, bool) {
	allFields := Fields{}
	merged := false

	for _, f := range fields {
		for key, value := range f {
			allFields[key] = value
			merged = true
		}
	}

	return allFields, merged
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
