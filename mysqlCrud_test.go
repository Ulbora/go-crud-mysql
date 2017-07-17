package mysqldb

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

var res bool
var testDb *sql.DB
var insertID int64

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
	var noTx *sql.DB
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
}

func TestDelete(t *testing.T) {
	var noTx *sql.DB
	var q = "DELETE FROM content WHERE id = ? "
	success := Delete(noTx, q, insertID)
	if success == true {
		fmt.Print("Deleted ")
		fmt.Println(insertID)
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
