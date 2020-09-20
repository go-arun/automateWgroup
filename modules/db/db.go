package db

import (
	"database/sql"
	"fmt"
	"strconv"
	// Is required to connect db
	_ "github.com/go-sql-driver/mysql"
	// "github.com/go-arun/fishrider/modules/session"
)
//Dbug enable this to ON debug mode , is no way realed to database module
var Dbug = false
//Connection .. to refer in all modules
var Connection *sql.DB

//Admin ... to pic admin details from db
type Admin struct {
	id int
	name,mail,sesid,pwd string
}

//Connect for connecting with MariaDB
func Connect(){
	dbDriver := "mysql"
	dbUser	:= "fradmin"
	dbPass	:= "f#$@rider"
	dbName	:= "db_fishrider"
	var err error
	Connection, err = sql.Open(dbDriver,dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println("dm-Error")
		panic(err.Error())
	}else if Dbug {//if in debug ON
		err = Connection.Ping() // Testing DB Connected or not 
		if err != nil {
			fmt.Println("dm-db.Ping failed:", err)
		}else {
			fmt.Println("dm-db.Ping Succecss:")
		}
	}
}

// item_desc varchar(45) 
// item_sel_price decimal(2,0) 
// item_buy_price decimal(2,0) 
// item_unit varchar(4) 
// item_stock int(11)
func InsertNewCatagory(itemName,itemUnit string)(string,error){
	insForm, err := Connection.Prepare(
		"INSERT INTO item_master(item_desc,item_unit) VALUES (?,?)",
	)
	_,err = insForm.Exec(itemName,itemUnit)
	if err != nil {
		return "",err // If insert failed return error Only chance is while trying for duplicate entry
	}

	//Geting last crated ID
	selDB, err := Connection.Query("SELECT max(item_id) FROM item_master")
	if err !=nil {
			panic(err.Error())
	}
	var itemID int
	for selDB.Next(){
		err = selDB.Scan(&itemID)
		if err != nil {
			panic(err.Error())
		}
	}
	
	return strconv.Itoa(itemID),nil
}