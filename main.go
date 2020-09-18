package main

import (
	"fmt"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func adminGet( c *gin.Context){
	if db.Dbug{
		fmt.Println("Inside func-adminGet")
	}
	c.HTML(
		http.StatusOK,
		"admin_login.html",
		gin.H{"title": "Admin Login"},
	)
}
func adminPost( c *gin.Context){
	if db.Dbug{
		fmt.Println("Inside func-adminPost")
	}
	c.HTML(
		http.StatusOK,
		"admin_login.html",
		gin.H{"title": "Admin Login"},
	)
}
func main(){
	dba := db.Connect()
	if db.Dbug {
		fmt.Println("%T",dba)
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", adminGet)
	router.POST("/admin", adminPost)
	router.Run()

}

