package sqldb

import "database/sql"

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/openTele?parseTime=true")
	if err != nil {
		panic(err)
	}

	return db
}
