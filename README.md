go-crud-mysql 
==============

Go Crud library for MySql

## TestInitialize

```
res = InitializeMysql("localhost:3306", "admin", "admin", "some_database")	
if res != true {
	// do something
}

```

## GetDb
### Gets a pointer to the DB for creating transactions

```
db = GetDb()
tx, err := db.Begin()

```

## Insert

```
var noTx *sql.Tx
var q = "INSERT INTO table (title, created_date, text, client_id) VALUES (?, ?, ?, ?)"
var a []interface{}
a = append(a, "test insert 2", time.Now(), "some content text", 125)	
success, insID := Insert(noTx, q, a...)

```

## Update

```
var noTx *sql.Tx
var q = "UPDATE table set title = ?, modified_date = ?, text = ? where id = ? and client_id = ? "
a := []interface{}{"test insert update", time.Now(), "some content new text updated", insertID, 125}
success := Update(noTx, q, a...)

```

## Get

```
a := []interface{}{insertID, 125}
var q = "select * from table WHERE id = ? and client_id = ?"
rowPtr := Get(q, a...)
if rowPtr != nil {
	foundRow := rowPtr.row
}

```

## GetList

```
a := []interface{}{125}
var q = "select * from table WHERE client_id = ? order by id"
rowsPtr := GetList(q, a...)
if rowsPtr != nil {	
	foundRows := rowsPtr.rows
}

```

## Delete

```
var noTx *sql.Tx
var q = "DELETE FROM table WHERE id = ? "
success := Delete(noTx, q, insertID)

```