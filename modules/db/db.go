package db

import (
	"database/sql"
	"fmt"
	// Is required to connect db
	_ "github.com/go-sql-driver/mysql"
)
//Dbug enable this to ON debug mode , is no way realed to database module
var Dbug = false

//Connect for connecting with MariaDB
func Connect() (db *sql.DB){
	dbDriver := "mysql"
	dbUser	:= "fradmin"
	dbPass	:= "f#$@rider"
	dbName	:= "db_fishrider"
	db, err := sql.Open(dbDriver,dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println("Error")
		panic(err.Error())
	}else if Dbug {//if in debug ON
		err = db.Ping() // Testing DB Connected or not 
		if err != nil {
			fmt.Println("db.Ping failed:", err)
		}else {
			fmt.Println("db.Ping Succecss:")
		}
	}
	return db
}
