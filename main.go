package main

import (
	"fmt"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/go-arun/fishrider/modules/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"log"
	"io"

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
func adminPanelPost(c *gin.Context){
	file, header, err := c.Request.FormFile("file") 
	if err != nil {
	  c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
	  return
	}
	filename := header.Filename
	out, err := os.Create("public/" + filename)
	if err != nil {
	  log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
	  log.Fatal(err)
	}
 
}





const (
	SR_File_Max_Bytes = 1024 * 1024 * 2
)

func main(){
	db.Connect() //db Connection 
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", adminGet)
	router.POST("/admin", adminPost)
	router.POST("/admin_panel", adminPanelPost)
	router.StaticFS("/file", http.Dir("public"))


	router.Run()
}

