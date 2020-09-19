package db

import (
	"database/sql"
	"fmt"
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
