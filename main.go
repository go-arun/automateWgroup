package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/go-arun/fishrider/modules/fileoperation"
	"github.com/go-arun/fishrider/modules/session"
)

//Item ... to store item details while retriving from DB
type Item struct {
	ID, Stock  int
	Name, Unit string
	Rate       float64
	IsUnitKg   bool // will display numbers/KG based on this, while listing items in user_index.html page
}

func adjStock(c *gin.Context) {
	itemID := c.PostForm("fish_name")
	actualStock := c.PostForm("adj_qty")
	strSQL := "UPDATE item_master SET item_stock = ? WHERE item_id = ?"
	updForm, err := db.Connection.Prepare(strSQL)
	if err != nil {
		panic(err.Error())
	}
	updForm.Exec(actualStock, itemID)
	//  c.Redirect(http.StatusTemporaryRedirect, "/admin_panel") // redirecting to admin loging page
	//  c.Abort()
	c.HTML(
		http.StatusOK,
		"admin_panel.html",
		gin.H{"title": "Admin Panel"},
	)

}
func updateStock(c *gin.Context) {
	itemID := c.PostForm("fishname")
	pPrice := c.PostForm("purch_price")
	sPrice := c.PostForm("selling_price")
	pQty := c.PostForm("purch_qty")
	strSQL := "UPDATE item_master SET item_stock = item_stock + ?,item_sel_price = ?,item_buy_price = ? WHERE item_id = ?"
	updForm, err := db.Connection.Prepare(strSQL)
	if err != nil {
		panic(err.Error())
	}

	updForm.Exec(pQty, sPrice, pPrice, itemID)
	c.HTML(
		http.StatusOK,
		"admin_panel.html",
		gin.H{"title": "Admin Panel"},
	)
}
func populateCategoryItems(c *gin.Context) {
	selDB, err := db.Connection.Query("SELECT item_id,item_desc,item_unit,item_stock FROM item_master")
	if err != nil {
		panic(err.Error())
	}

	item := Item{}
	itemCollection := []Item{}

	for selDB.Next() {
		var id int
		var name string
		var unit string
		var stk int
		err = selDB.Scan(&id, &name, &unit, &stk)
		if err != nil {
			panic(err.Error())
		}
		item.ID = id
		item.Name = name
		item.Unit = unit
		item.Stock = stk

		itemCollection = append(itemCollection, item)
	}
	c.HTML(
		http.StatusOK,
		"admin_panel.html", gin.H{
			"FishCatagories": itemCollection,
		})
}

func addNewCatagory(c *gin.Context) {
	category := c.PostForm("category")
	unit := c.PostForm("unit")
	file, header, err := c.Request.FormFile("filename")
	if db.Dbug {
		fmt.Println("New category and unit is --", category, unit)
	}
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	//Inserting New Catagory into DB
	itemID, err := db.InsertNewCatagory(category, unit)
	if err != nil { // TBD : If insert failed ,it is showing error diretly in page , user need an alert
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	filename = fileoperation.ReplaceFileName(filename, itemID) // Replace Files 'name' part with itemID
	if db.Dbug {
		fmt.Println("file Name Before Replacing & ID returned from DB --->:", filename, itemID)
	}
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

func adminGet(c *gin.Context) {
	//Checking for any active sessions
	IsSectionActive := session.SessinStatus(c, "admin_session_cookie")
	if IsSectionActive {
		populateCategoryItems(c)
	} else {
		fmt.Println("No Active Sessions found ")
		// c.HTML(http.StatusOK, "admin_login.html", []string{"a", "b", "c"})
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "Admin Login",
				"diplay": "none",
			},
		)
	}
}
func adminPost(c *gin.Context) {
	if db.Dbug {
		fmt.Println("dm-Inside func-adminPost")
	}
	c.Request.ParseForm()
	name := c.PostForm("uname")
	pwd := c.PostForm("pwd")
	// name := c.Request.PostForm["uname"][0]
	// pwd  := c.Request.PostForm["pwd"][0]

	if session.UsrCredentialsVerify(name, pwd) { //return value 'true' means creadentias are matching ..
		//SetNewSessionID
		session.SetSessionCookie(c, name, "admin_session_cookie")
		populateCategoryItems(c)
	} else {
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "Admin Login",
				"diplay": "block",
			},
		)
	}

}
func adminPanelPost(c *gin.Context) {
	operation := c.PostForm("action") // based on this , will decide which operation need to done
	fmt.Println("User Opted action is :  ", operation)
	switch operation {
	//add new catagory
	case "category":
		addNewCatagory(c)
	case "purchase":
		updateStock(c)
	case "stockadj":
		adjStock(c)

	}

}

func adminPanelGet(c *gin.Context) {
	//Checking for any active sessions
	IsSectionActive := session.SessinStatus(c, "admin_session_cookie")
	if IsSectionActive {
		populateCategoryItems(c)
	} else {
		fmt.Println("No Active Sessions found ")
		c.Redirect(http.StatusTemporaryRedirect, "/admin") // redirecting to admin loging page
		c.Abort()
	}
}

func logoutGet(c *gin.Context) {
	session.RemoveSessionIDFromDB(c)
	c.Redirect(http.StatusTemporaryRedirect, "/admin") // redirecting to admin loging page
	c.Abort()
}

func userIndexGet(c *gin.Context) {
	selDB, err := db.Connection.Query("SELECT item_id,item_desc,item_unit,item_stock,item_sel_price FROM item_master WHERE item_stock > 0")
	if err != nil {
		panic(err.Error())
	}

	item := Item{}
	itemCollection := []Item{}

	for selDB.Next() {
		item.IsUnitKg = false // just resetting back to fales ( TBD : make sure this is necessary)
		var id int
		var name string
		var unit string
		var stk int
		var price float64
		err = selDB.Scan(&id, &name, &unit, &stk, &price)
		if err != nil {
			panic(err.Error())
		}
		if unit == "Kg" {
			item.IsUnitKg = true
		}
		item.ID = id
		item.Name = name
		item.Unit = unit
		item.Stock = stk
		item.Rate = price
		itemCollection = append(itemCollection, item)
	}
	c.HTML(
		http.StatusOK,
		"user_index.html", gin.H{
			"AvailableSock": itemCollection,
		})

}
func main() {
	db.Connect() //db Connection
	router := gin.Default()
	//admin routes
	router.LoadHTMLGlob("templates/*")
	router.GET("/admin", adminGet)
	router.POST("/admin", adminPost)
	router.POST("/admin_panel", adminPanelPost)
	router.GET("/admin_panel", adminPanelGet)
	router.GET("/logout", logoutGet)
	router.StaticFS("/file", http.Dir("pics"))
	//user Routes
	router.GET("/", userIndexGet)
	router.Run()
}
