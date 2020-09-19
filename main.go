package main

import (
	"fmt"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/go-arun/fishrider/modules/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

func adminGet( c *gin.Context){
	if db.Dbug{
		fmt.Println("dm-Inside func-adminGet")
	}
	c.HTML(
		http.StatusOK,
		"admin_login.html",
		gin.H{"title": "Admin Login"},
	)
}

func adminPost( c *gin.Context){
	if db.Dbug{fmt.Println("dm-Inside func-adminPost")}
	c.Request.ParseForm()
	name := c.Request.PostForm["uname"][0]
	pwd  := c.Request.PostForm["pwd"][0]
	session.UsrCredentialsVerify(name,pwd)
	c.HTML(
		http.StatusOK,
		"admin_panel.html",
		gin.H{"title": "Admin Panel"},
	)
}
func main(){
	db.Connect() //db Connection 
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", adminGet)
	router.POST("/admin", adminPost)
	router.Run()
}

