package mysqldb

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var res bool
var testDb *sql.DB
var insertID int64
var insertID2 int64
var insertID3 int64
var insertID4 int64

func TestInitialize(t *testing.T) {
	res = InitializeMysql("localhost:3306", "admin", "admin", "ulbora_content_service")
	//fmt.Print("res in init: ")
	//fmt.Println(res)
	if res != true {
		fmt.Println("database init failed")
		t.Fail()
	}
}

func TestGetDb(t *testing.T) {
	testDb = GetDb()
	if db == nil {
		fmt.Println("get db failed")
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	var noTx *sql.Tx
	var q = "INSERT INTO content (title, created_date, text, client_id) VALUES (?, ?, ?, ?)"
	var a []interface{}
	a = append(a, "test insert 2", time.Now(), "some content text", 125)
	//can also be: a := []interface{}{"test insert", time.Now(), "some content text", 125}
	success, insID := Insert(noTx, q, a...)
	if success == true && insID != -1 {
		insertID = insID
		fmt.Print("new Id: ")
		fmt.Println(insID)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}

	success, insID2 := Insert(noTx, q, a...)
	if success == true && insID2 != -1 {
		insertID2 = insID2
		fmt.Print("new Id: ")
		fmt.Println(insID2)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	var noTx *sql.Tx
	var q = "UPDATE content set title = ?, modified_date = ?, text = ? where id = ? and client_id = ? "
	//var a []interface{}
	//a = append(a, "test insert 2", time.Now(), "some content text", 125)
	a := []interface{}{"test insert update", time.Now(), "some content new text updated", insertID, 125}
	success := Update(noTx, q, a...)
	if success == true {
		fmt.Println("database update success")
	} else {
		fmt.Println("database update failed")
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	a := []interface{}{insertID, 125}
	var q = "select * from content WHERE id = ? and client_id = ?"
	rowPtr := Get(q, a...)
	if rowPtr != nil {
		//fmt.Print("columns")
		//fmt.Println(rowPtr.columns)
		foundRow := rowPtr.Row
		//fmt.Print("Get ")
		//fmt.Println(foundRow)
		//fmt.Println("Get results: --------------------------")
		int64Val, err2 := strconv.ParseInt(foundRow[0], 10, 0)
		if err2 != nil {
			fmt.Print(err2)
		}
		if insertID != int64Val {
			fmt.Print(insertID)
			fmt.Print(" != ")
			fmt.Println(int64Val)
			t.Fail()
		}
	} else {
		fmt.Println("database getRow failed")
		t.Fail()
	}
}

func TestGetList(t *testing.T) {
	a := []interface{}{125}
	var q = "select * from content WHERE client_id = ? order by id"
	rowsPtr := GetList(q, a...)
	if rowsPtr != nil {
		//fmt.Print("columns")
		//fmt.Println(rowsPtr.columns)
		foundRows := rowsPtr.Rows
		//fmt.Print("GetList ")
		//fmt.Println(foundRows)
		//fmt.Println("GetList results: --------------------------")
		for r := range foundRows {
			foundRow := foundRows[r]
			for c := range foundRow {
				if c == 0 {
					int64Val, err2 := strconv.ParseInt(foundRow[c], 10, 0)
					if err2 != nil {
						fmt.Print(err2)
					}
					if r == 0 {
						if insertID != int64Val {
							fmt.Print(insertID)
							fmt.Print(" != ")
							fmt.Println(int64Val)
							t.Fail()
						}
					} else if r == 1 {
						if insertID2 != int64Val {
							fmt.Print(insertID)
							fmt.Print(" != ")
							fmt.Println(int64Val)
							t.Fail()
						}
					}
				}
				//fmt.Println(string(foundRow[c]))
			}
		}
	} else {
		fmt.Println("database getRow failed")
		t.Fail()
	}
}
func TestDelete(t *testing.T) {
	a := []interface{}{insertID}
	a2 := []interface{}{insertID2}
	var noTx *sql.Tx
	var q = "DELETE FROM content WHERE id = ? "
	success := Delete(noTx, q, a...)
	if success == true {
		fmt.Print("Deleted ")
		fmt.Println(insertID)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}

	success2 := Delete(noTx, q, a2...)
	if success2 == true {
		fmt.Print("Deleted ")
		fmt.Println(insertID2)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}
}

func TestInsertTx(t *testing.T) {
	db := GetDb()
	var success = false
	var q = "INSERT INTO content (title, created_date, text, client_id) VALUES (?, ?, ?, ?)"
	var a []interface{}
	a = append(a, "test insert with tx", time.Now(), "some content text", 125)
	//can also be: a := []interface{}{"test insert", time.Now(), "some content text", 125}
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	success1, insID1 := Insert(tx, q, a...)
	success2, insID2 := Insert(tx, q, a...)
	defer func() {
		if success != true {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if success1 == true && insID1 != -1 && success2 == true && insID2 != -1 {
		insertID3 = insID1
		insertID4 = insID2
		success = true
		//fmt.Print("new Id with tx: ")
		//fmt.Println(insID)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}

}

func TestDeleteTx(t *testing.T) {
	a3 := []interface{}{insertID3}
	a4 := []interface{}{insertID4}
	var noTx *sql.Tx
	var q = "DELETE FROM content WHERE id = ? "
	success := Delete(noTx, q, a3...)
	if success == true {
		fmt.Print("Deleted ")
		fmt.Println(insertID3)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}

	success2 := Delete(noTx, q, a4...)
	if success2 == true {
		fmt.Print("Deleted ")
		fmt.Println(insertID4)
	} else {
		fmt.Println("database insert failed")
		t.Fail()
	}
}
func TestClose(t *testing.T) {
	if res == true {
		rtn := Close()
		if rtn != true {
			fmt.Println("database close failed")
			t.Fail()
		}
	}
}
