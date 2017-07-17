package mysqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

//InitializeMysql Mysql init to mysql
func InitializeMysql(host, user, pw, dbName string) bool {
	var rtn = false

	var conStr = user + ":" + pw + "@tcp(" + host + ")/" + dbName
	//var conStr = user + ":" + pw + "@/" + dbName
	//fmt.Println("conStr: " + conStr)

	//db, err := sql.Open("mysql", "user:password@/dbname")
	db, err = sql.Open("mysql", conStr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	} else {
		rtn = true
	}
	return rtn

}

//GetDb gets a pointer to db for transactions
func GetDb() *sql.DB {
	return db
}

//Insert inserts a row. Passing in tx allows for transactions
func Insert(tx *sql.DB, query string, args ...interface{}) (bool, int64) {
	var success = false
	var id int64 = -1
	var dbToUse *sql.DB
	if tx != nil {
		dbToUse = tx
	} else {
		dbToUse = db
	}
	stmtIns, err := dbToUse.Prepare(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()
	//res, err := stmtIns.Exec("test", time.Now(), "test", 111)
	res, err := stmtIns.Exec(args...)
	if err != nil {
		fmt.Println("Insert Exec err:", err.Error())
	} else {
		fmt.Println("Insert Exec success:")
		id, err = res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		} else {
			//println("LastInsertId:", id)
			success = true
		}
	}
	return success, id
}

//Delete deletes records
func Delete(tx *sql.DB, query string, id int64) bool {
	var success = false
	var dbToUse *sql.DB
	if tx != nil {
		dbToUse = tx
	} else {
		dbToUse = db
	}
	stmtIns, err := dbToUse.Prepare(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()
	//res, err := stmtIns.Exec("test", time.Now(), "test", 111)
	res, err := stmtIns.Exec(id)
	if err != nil {
		fmt.Println("Delete Exec err:", err.Error())
	} else {
		affect, err := res.RowsAffected()
		if err != nil {
			println("Error:", err.Error())
		} else {
			fmt.Println(affect)
			fmt.Println("Delete Exec success:")
			success = true
		}
	}
	return success
}

//Close close
func Close() bool {
	var rtn = false
	db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		rtn = true
	}
	return rtn
}
