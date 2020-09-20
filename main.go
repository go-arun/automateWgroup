package main

import (
	"fmt"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/go-arun/fishrider/modules/session"
	"github.com/go-arun/fishrider/modules/fileoperation"
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
	category := c.PostForm("category")
	unit  := c.PostForm("unit")
	file, header, err := c.Request.FormFile("filename") 
	if db.Dbug{fmt.Println("New category and unit is --",category,unit)}

	//Inserting New Catagory into DB
	itemID := db.InsertNewCatagory(category,unit)
	if err != nil {
	  c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
	  return
	}
	filename := header.Filename
	fmt.Println("file Name Before Replacing & ID returned from DB --->:",filename,itemID)
	filename = fileoperation.ReplaceFileName(filename,itemID) // Replace Files 'name' part with itemID
	fmt.Println("file Name After Replacing  --->:",filename)

	fmt.Println("FileName Changed to ",filename)
	out, err := os.Create("pics/" + filename)
	if err != nil {
	  log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
	  log.Fatal(err)
	}
 }

func main(){
	db.Connect() //db Connection 
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", adminGet)
	router.POST("/admin", adminPost)
	router.POST("/admin_panel", adminPanelPost)
	router.StaticFS("/file", http.Dir("pics"))
	router.Run()
}

