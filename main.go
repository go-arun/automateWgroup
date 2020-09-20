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
//Item ... to store item details while retriving from DB
type Item struct {
	ID int
	Name,Unit string
}
func populateCategoryItems(c *gin.Context){
	selDB, err := db.Connection.Query("SELECT item_id,item_desc,item_unit FROM item_master")
	if err !=nil {
			panic(err.Error())
	}

	item := Item{}
	itemCollection := []Item{}

	for selDB.Next(){
		var id int
		var name string
		var unit string
		err = selDB.Scan(&id,&name,&unit)
		if err != nil {
			panic(err.Error())
		}
		item.ID = id
		item.Name = name
		item.Unit = unit
		itemCollection = append(itemCollection,item)
	}
	c.HTML(											
		http.StatusOK,
		"admin_panel.html",gin.H{
			"FishCatagories": itemCollection,
		})
}

func addNewCatagory(c *gin.Context){
	category := c.PostForm("category")
	unit  := c.PostForm("unit")
	file, header, err := c.Request.FormFile("filename") 
	if db.Dbug{fmt.Println("New category and unit is --",category,unit)}
	if err != nil {
	  c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
	  return
	}
	//Inserting New Catagory into DB
	itemID,err := db.InsertNewCatagory(category,unit)
	if err != nil{// TBD : If insert failed ,it is showing error diretly in page , user need an alert 
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	filename = fileoperation.ReplaceFileName(filename,itemID) // Replace Files 'name' part with itemID
	if db.Dbug {fmt.Println("file Name Before Replacing & ID returned from DB --->:",filename,itemID)}
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
	populateCategoryItems(c)

}
func adminPanelPost(c *gin.Context){
	operation :=  c.PostForm("action") // baed on this , will decide which operation need to done 
	switch operation {
	//add new catagory	
	case "category":
		addNewCatagory(c)
	case "purchase":
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

