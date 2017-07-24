package mysqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

//DbRow database row
type DbRow struct {
	columns []string
	row     []string
}

//DbRows array of database rows
type DbRows struct {
	columns []string
	//rows    [][]sql.RawBytes
	rows [][]string
}

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
func Insert(tx *sql.Tx, query string, args ...interface{}) (bool, int64) {
	var success = false
	var id int64 = -1
	var stmtIns *sql.Stmt
	if tx != nil {
		fmt.Println("Using tx")
		stmtIns, err = tx.Prepare(query)
	} else {
		stmtIns, err = db.Prepare(query)
	}

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
			fmt.Println("Error:", err.Error())
		} else {
			//fmt.Println("LastInsertId:", id)
			success = true
		}
	}

	return success, id
}

//Update updates a row. Passing in tx allows for transactions
func Update(tx *sql.DB, query string, args ...interface{}) bool {
	var success = false
	var dbToUse *sql.DB
	if tx != nil {
		dbToUse = tx
	} else {
		dbToUse = db
	}
	stmtUp, err := dbToUse.Prepare(query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUp.Close()
	//res, err := stmtIns.Exec("test", time.Now(), "test", 111)
	res, err := stmtUp.Exec(args...)
	if err != nil {
		fmt.Println("Update Exec err:", err.Error())
	} else {
		fmt.Println("Update Exec success:")
		affectedRows, err := res.RowsAffected()
		//fmt.Println(affectedRows)
		if err != nil && affectedRows == 0 {
			fmt.Println("Error:", err.Error())
		} else {
			//fmt.Println("LastInsertId:", id)
			success = true
		}
	}
	return success
}

//Get get a row. Passing in tx allows for transactions
func Get(query string, args ...interface{}) *DbRow {
	var rtn DbRow
	stmtGet, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmtGet.Close()
	rows, err := stmtGet.Query(args...)

	defer rows.Close()
	if err != nil {
		fmt.Print("Get err: ")
		fmt.Println(err)
	} else {
		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error())
		}
		rtn.columns = columns
		rowValues := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(rowValues))
		for i := range rowValues {
			scanArgs[i] = &rowValues[i]
		}
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err.Error())
			}
			for _, col := range rowValues {
				var value string
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				//fmt.Println(columns[i], ": ", value)
				rtn.row = append(rtn.row, value)
			}
			//rtn.row = rowValues
		}
		//fmt.Print("rtn.row")
		//fmt.Println(rtn.row)
		if err = rows.Err(); err != nil {
			panic(err.Error())
		}
	}
	return &rtn
}

//GetList get a list of rows. Passing in tx allows for transactions
func GetList(query string, args ...interface{}) *DbRows {
	var rtn DbRows
	stmtGet, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmtGet.Close()
	rows, err := stmtGet.Query(args...)
	defer rows.Close()
	if err != nil {
		fmt.Print("GetList err: ")
		fmt.Println(err)
	} else {
		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error())
		}
		rtn.columns = columns
		rowValues := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(rowValues))
		for i := range rowValues {
			scanArgs[i] = &rowValues[i]
		}
		for rows.Next() {
			var rowValuesStr []string
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err.Error())
			}
			//fmt.Print("len(rowValues:)")
			//fmt.Println(len(rowValues))

			for _, col := range rowValues {
				var value string
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				//fmt.Println(columns[i], ": ", value)
				rowValuesStr = append(rowValuesStr, value)
			}
			//fmt.Print("rowValuesStr: ")
			//fmt.Println(rowValuesStr)
			rtn.rows = append(rtn.rows, rowValuesStr)
		}
		//fmt.Print("rtn.rows:")
		//fmt.Println(rtn.rows)
		if err = rows.Err(); err != nil {
			panic(err.Error())
		}
	}
	return &rtn
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
		affectedRows, err := res.RowsAffected()
		if err != nil && affectedRows == 0 {
			fmt.Println("Error:", err.Error())
		} else {
			//fmt.Println(affectedRows)
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
