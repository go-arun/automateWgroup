package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/go-arun/fishrider/modules/fileoperation"
	"github.com/go-arun/fishrider/modules/session"
	"github.com/go-arun/fishrider/modules/smsapi"
	razorpay "github.com/razorpay/razorpay-go"
)

// to store otp_id received from SMS API Provider
var otpIDfromProvider string

// to store user mobile number
var globalMobNo int

//Item ... to store item details while retriving from DB
type Item struct {
	ID, Stock  int
	Name, Unit string
	Rate       float64
	PurchRate  float64
	IsUnitKg   bool // will display numbers/KG based on this, while listing items in user_index.html page

}

//CartItem ... to store cart items retrieved from Cokkies to Render
type CartItem struct {
	SlNo int
	Desc string
	Qty  int
	Rate string
	// Rate     float64
	Unit     string
	SubTotal string
	// SubTotal float64
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
func populateCategoryItems(c *gin.Context, itemID string) { // later move to db Module TBD
	var strSQL string // that time remove 2ndArgument, may be useless that time TBD
	if itemID == "" { //Caller need whole details
		strSQL = "SELECT item_id,item_desc,item_unit,item_stock,item_sel_price,item_buy_price FROM item_master ORDER BY  item_stock DESC"
	} else { // choose only one item based on given itemID
		strSQL = "SELECT item_id,item_desc,item_unit,item_stock,item_sel_price,item_buy_price FROM item_master WHERE item_id=" + itemID
	}
	selDB, err := db.Connection.Query(strSQL)
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
		var itemSelPrice, itemBuyprice float64
		err = selDB.Scan(&id, &name, &unit, &stk, &itemSelPrice, &itemBuyprice)
		if err != nil {
			panic(err.Error())
		}
		item.ID = id
		item.Name = name
		item.Unit = unit
		item.Stock = stk
		item.Rate = itemSelPrice
		item.PurchRate = itemBuyprice

		itemCollection = append(itemCollection, item)
	}
	if itemID == "" { // Means taken all details to show in stock list option
		c.HTML(
			http.StatusOK,
			"admin_panel.html", gin.H{
				"FishCatagories": itemCollection,
			})
	} else { // Just picked asked item details only for view/edit
		c.HTML(
			http.StatusOK,
			"edit_item.html", gin.H{
				"FishCatagories": itemCollection,
				"diplay":         "none",
			})
	}

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
	populateCategoryItems(c, "") // Everything fine in adding new catagory , but I am not happy , req err handling TBD
}

func adminGet(c *gin.Context) {
	//Checking for any active sessions
	IsSectionActive, _ := session.SessinStatus(c, "admin_session_cookie")
	if IsSectionActive {
		populateCategoryItems(c, "")
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

	if session.AdminCredentialsVerify(name, pwd) { //return value 'true' means creadentias are matching ..
		//SetNewSessionID
		session.SetAdminSessionCookie(c, name, "admin_session_cookie")
		populateCategoryItems(c, "")
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
	IsSectionActive, _ := session.SessinStatus(c, "admin_session_cookie")
	if IsSectionActive {
		populateCategoryItems(c, "")
	} else {
		fmt.Println("No Active Sessions found ")
		c.Redirect(http.StatusTemporaryRedirect, "/admin") // redirecting to admin loging page
		c.Abort()
	}
}

func logoutGet(c *gin.Context) {
	session.RemoveAdminSessionIDFromDB(c)
	c.Redirect(http.StatusTemporaryRedirect, "/admin") // redirecting to admin loging page
	c.Abort()
}

func userIndexGet(c *gin.Context) {
	// Logic to show User name in navbar
	IsUsrSectionActive, usrName := session.SessinStatus(c, "user_session_cookie")
	if !IsUsrSectionActive {
		usrName = "Sign In"
	}

	//For safer side alway clear cart
	session.Cart = nil

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
			"usrName":       usrName,
		})

}
func userSignupGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"user_signup.html",
		gin.H{"title": "User SignUp",
			"diplay": "none", // TBD make use of this logic to diplay error
		},
	)

}
func userSignupPost(c *gin.Context) {
	mob := c.PostForm("mobile")
	name := c.PostForm("cust_name")
	houseName := c.PostForm("house_name")
	street := c.PostForm("street_name")
	landMark := c.PostForm("lmark_name")
	fmt.Println(mob, name, houseName, street, landMark)
	err := db.AddNewCustomer(mob, name, houseName, street, landMark)
	if err != nil {
		fmt.Println("ERROR inserting new user")
	}

	//After signup continue to Login Page
	c.HTML(
		http.StatusOK,
		"user_login.html", gin.H{
			"title": "User Registration",
		})

}

//items selected and moving to Orders page
func userIndexPost(c *gin.Context) {
	lastItmID := db.GetLastItemID() // iterating from 0 to lastItmID val, to track
	// the selcted items by user.( html select id names are itemIds)
	for lastItmID > 0 { // Storing selected items as Cookie
		if c.PostForm(strconv.Itoa(lastItmID)) != "Qty" && c.PostForm(strconv.Itoa(lastItmID)) != "" {
			if db.Dbug {
				fmt.Println("Items Selected-:", c.PostForm(strconv.Itoa(lastItmID)))
			}
			session.PushSelectionToCookie(c, strconv.Itoa(lastItmID), c.PostForm(strconv.Itoa(lastItmID)))
		}
		lastItmID = lastItmID - 1
	}

	//Checking for any active sessions
	IsUsrSectionActive, _ := session.SessinStatus(c, "user_session_cookie")
	if IsUsrSectionActive { // Moving to Orders page
		fmt.Println("user asession is actice -- So Moving from landing page to Oreder ")
		c.Redirect(http.StatusTemporaryRedirect, "/orders")
		c.Abort()
		// c.HTML(
		// 	http.StatusOK,
		// 	"orders.html",
		// 	gin.H{"title": "User Login",
		// 		"diplay": "none", // TBD make use of this logic to diplay error
		// 	},
		// )
		fmt.Println("Session is active for this user")
	} else {
		fmt.Println("No Active USR Sessions found ")
		//Login User
		c.HTML(
			http.StatusOK,
			"user_login.html",
			gin.H{"title": "User Login",
				"diplay": "none", // TBD make use of this logic to diplay error
			},
		)
	}

}

func userLoginGet(c *gin.Context) {
	//Checking for any active sessions
	IsUsrSectionActive, _ := session.SessinStatus(c, "user_session_cookie")
	if IsUsrSectionActive {
		//Move to Item lising page TBD
	} else {
		fmt.Println("No Active USR Sessions found ")
		// c.HTML(http.StatusOK, "admin_login.html", []string{"a", "b", "c"})
		c.HTML(
			http.StatusOK,
			"user_login.html",
			gin.H{"title": "User Login",
				"diplay": "none", // TBD make use of this logic to diplay error
			},
		)
	}

}
func userLoginPost(c *gin.Context) {

	userMob := c.PostForm("mobile")
	RegisteredUser, _, _ := session.UserCredentialsVerify(userMob)
	if !RegisteredUser { // This mobile not in our record , Direct to register page
		c.HTML(
			http.StatusOK,
			"user_signup.html",
			gin.H{"title": "User SignUp",
				"diplay": "block", // TBD make use of this logic to diplay error
			},
		)

	} else {
		otpIDfromProvider = smsapi.GenerateOTP(userMob)
		globalMobNo, _ = strconv.Atoi(userMob)
		c.HTML(
			http.StatusOK,
			"user_otp_verification.html",
			gin.H{"title": "OTP Verification",
				"diplay": "none", // TBD make use of this logic to diplay error
			},
		)
	}
}
func userOtpVerifyPost(c *gin.Context) {
	codeP1 := c.PostForm("code_p1")
	codeP2 := c.PostForm("code_p2")
	codeP3 := c.PostForm("code_p3")
	codeP4 := c.PostForm("code_p4")
	codeP5 := c.PostForm("code_p5")
	codeP6 := c.PostForm("code_p6")
	wholeCode := codeP1 + codeP2 + codeP3 + codeP4 + codeP5 + codeP6
	mobileVerified, _ := smsapi.VerifyOTP(otpIDfromProvider, wholeCode)
	if mobileVerified {
		session.SetUserSessionCookie(c, globalMobNo, "user_session_cookie")
		fmt.Println("OTP Verification Success")
	} else {
		fmt.Println("OTP Verification Success")
	}

	//User Logged so Move to Orders Page
	c.Redirect(http.StatusTemporaryRedirect, "/orders")
	c.Abort()
	// c.HTML(
	// 	http.StatusOK,
	// 	"orders.html",
	// 	gin.H{"title": "User Login",
	// 		"diplay": "none", // TBD make use of this logic to diplay error
	// 	},
	// )

}
func userLogoutGet(c *gin.Context) {
	session.RemoveUserSessionIDFromDB(c)
	c.Redirect(http.StatusTemporaryRedirect, "/") // redirecting to item listing page
	c.Abort()

}

func userOrdersGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"orders.html",
		gin.H{"title": "User Login",
			"diplay": "none", // TBD make use of this logic to diplay error
		},
	)
}
func userOrdersPost(c *gin.Context) { // Redirected from IndexpagePost event
	// Logic to show User name in navbar
	IsUsrSectionActive, usrName := session.SessinStatus(c, "user_session_cookie")
	if !IsUsrSectionActive {
		usrName = "Sign In"
	}
	cartItems := session.PullCartItemFromCookie(c) // Return a struct array of cart items retrived from Cookies
	var singleCartItem CartItem
	var fullCartItems []CartItem
	var TotalAmt float64
	for key := range cartItems { // range through the array contains the cookie(havning only icode and qty) and adding missing details from DB
		singleCartItem.SlNo = key + 1
		singleCartItem.Qty, _ = strconv.Atoi(cartItems[key].IQty)
		desc, rate, unit, _, _, _ := db.GetItemDetails(cartItems[key].ICode)
		singleCartItem.Desc = desc
		singleCartItem.Rate = fmt.Sprintf("%.2f", rate)
		singleCartItem.Unit = unit
		singleCartItem.SubTotal = fmt.Sprintf("%.2f", float64(singleCartItem.Qty)*rate)

		fullCartItems = append(fullCartItems, singleCartItem)
		fmt.Println(cartItems[key].ICode)
		//fmt.Println("key and val=", key, val)
		subTotalToFloat, _ := strconv.ParseFloat(singleCartItem.SubTotal, 64)
		TotalAmt = TotalAmt + subTotalToFloat
	}
	TotalAmtInPaisa := TotalAmt * 100 // This is required while initate for payment in Razorpay

	TotalAmtString := fmt.Sprintf("%.2f", TotalAmt)
	c.HTML(
		http.StatusOK,
		"orders.html",
		gin.H{"title": "User Login",
			"ItemsOrdered":    fullCartItems,
			"TotalAmt":        TotalAmtString,
			"TotalAmtInPaisa": TotalAmtInPaisa,
			"usrName":         usrName,
		},
	)

}

func orderconfirmPost(c *gin.Context) { // Execuet this after selecting payment mode ( Confirm Order)
	sessionCookie, _ := c.Cookie("user_session_cookie")
	_, cMo, cNme, cAdr1, cAdr2, _ := db.TraceUserWithSIDinDB(sessionCookie)
	amtInPaisa, _ := strconv.Atoi(c.PostForm("amt_inPaisa")) //Eg 24545
	totalamt := c.PostForm("orderAmt")                       // in noraml foramt Eg 245.45
	if c.PostForm("paymentMode") == "online" {
		client := razorpay.NewClient("rzp_test_zlsYsrvuUxxhln", "9LtUy4qpLCOtl4Gz2asp59es")
		data := map[string]interface{}{
			"amount":   amtInPaisa,
			"currency": "INR",
			"receipt":  "110",
			// "receipt":      "some_receipt_id",
		}
		//var emptyMap = map[string]string{}
		hsh := make(map[string]string)
		body, err := client.Order.Create(data, hsh)
		if err != nil {
			fmt.Println(err)
		}
		if db.Dbug {
			fmt.Println("body RazorPay Response with order_id :", body)
		}
		RPayOrderID := body["id"]
		fmt.Println("RPayOrderID--->:", RPayOrderID)
		c.HTML(
			http.StatusOK,
			"payment.html",
			gin.H{"title": "User Login",
				"RPayOrderID": RPayOrderID,
				"amtInPaisa":  amtInPaisa,
				"UserName":    cNme,
				"Mobile":      cMo,
				"DelAddr":     cAdr1 + cAdr2,
				"TotalAmt":    totalamt,
			},
		)

	} else { // Cod Payemnt so will directy go to Orders Page and procedd next steps
		c.Redirect(http.StatusTemporaryRedirect, "/orderhistory") // redirecting
		c.Abort()
	}
}

//To show order details-while accessing from NAV bar
// if the user is logged in - otherwise skip to login page
func orderHistoryGet(c *gin.Context) {
	IsUsrSectionActive, usrName := session.SessinStatus(c, "user_session_cookie")
	if !IsUsrSectionActive {
		c.HTML(
			http.StatusOK,
			"user_login.html",
			gin.H{"title": "User Login",
				"diplay": "none", // TBD make use of this logic to diplay error
			},
		)

	} else {
		sessionCookie, _ := c.Cookie("user_session_cookie")
		_, _, _, _, _, custID := db.TraceUserWithSIDinDB(sessionCookie)
		//Collecting all Order details to show
		oK, UserOrderHisory := db.GetOrderHistory(custID)
		if !oK {
			fmt.Println("Something is went wrong while collecting order details !!")
		}
		c.HTML(
			http.StatusOK,
			"orderhistory.html",
			gin.H{"title": "Orders",
				"OrderHistrory": UserOrderHisory,
				"usrName":       usrName,
			},
		)

	}

}

func orderHistoryPost(c *gin.Context) { //Will execute this After placing an order(& payment)
	IsUsrSectionActive, usrName := session.SessinStatus(c, "user_session_cookie")
	if IsUsrSectionActive { // Dont want to entertain guest users herer !!!
		payMode := c.PostForm("paymentMode")
		if payMode != "cod" {
			payMode = "online"
		} // will not get the value if it comes from Online payment page, so setting it here
		//Insert Current Order details to to DB
		sessionCookie, _ := c.Cookie("user_session_cookie")
		_, _, _, _, _, custID := db.TraceUserWithSIDinDB(sessionCookie)
		toatlInString := c.PostForm("orderAmt") //Fun part is, in COD and ONLINE getting this value from 2 different forms , both input having the same name
		fmt.Println("toatlInString Vql:->", toatlInString)
		totalToFloat, err := strconv.ParseFloat(toatlInString, 64)
		if err != nil {
			fmt.Println("Convertion errror: ", err)
		}
		// This is only updating master order table
		_, newOrderID := db.AddNewOrderEntry(custID, totalToFloat, payMode)

		//Now add induvidual item details to order_details table
		cartItems := session.PullCartItemFromCookie(c)
		for key := range cartItems {
			iCode, _ := strconv.Atoi(cartItems[key].ICode)
			iQty, _ := strconv.Atoi(cartItems[key].IQty)
			fmt.Println("icode and iqty", iCode, iQty)
			okay := db.UpdateOrderDetails(newOrderID, iCode, iQty)
			if !okay {
				fmt.Println("Error in inserting to Order details ....")
			}
		}
		//Now safe to remove cart entries in Cookies
		session.RemoveCookie(c, "user_cart")
		cartItems = nil
		session.Cart = nil

		//Collecting all Order details to show
		oK, UserOrderHosory := db.GetOrderHistory(custID)
		if !oK {
			fmt.Println("Something is went wrong while collecting order details !!")
		}
		c.HTML(
			http.StatusOK,
			"orderhistory.html",
			gin.H{"title": "User Login",
				"OrderHistrory": UserOrderHosory,
				"usrName":       usrName,
			},
		)
	} else { // User is not logged in so do that first
		c.HTML(
			http.StatusOK,
			"user_login.html",
			gin.H{"title": "User Login",
				"diplay": "none", // TBD make use of this logic to diplay error
			},
		)
	}

}

// To disply any particulr Order
func viewAnyOrderGet(c *gin.Context) {
	OrdID := c.Request.URL.Query()["ordid"][0] // Getting Order ID passed with URL
	_, usrName := session.SessinStatus(c, "user_session_cookie")
	fmt.Println("Wnat to see the order details of order number ", OrdID)
	oK, itemsList, date, status, PayMode, amt := db.GetSingleOredrDetails(OrdID)
	if !oK {
		fmt.Println("Something went wrong while picking Single Order Deatils ..Please have a look")
	}
	fmt.Println(oK, itemsList, date, status, PayMode, amt)
	//		subTotalToFloat, _ := strconv.ParseFloat(singleCartItem.SubTotal, 64)
	//		TotalAmt = TotalAmt + subTotalToFloat
	//	TotalAmtInPaisa := TotalAmt * 100 // This is required while initate for payment in Razorpay

	//	TotalAmtString := fmt.Sprintf("%.2f", TotalAmt)

	c.HTML(
		http.StatusOK,
		"view_particular_order.html",
		gin.H{"title": "OrderDetail",
			"ItemsOrdered": itemsList,
			"OrdID":        OrdID,
			"date":         date,
			"PayMode":      PayMode,
			"amt":          amt,
			"OrdStatus":    status,
			"usrName":      usrName,

			// "TotalAmt":        TotalAmtString,
			// "TotalAmtInPaisa": TotalAmtInPaisa,
		},
	)

}

//Edit any sigle item from Slock List
func viewandedititemGet(c *gin.Context) {
	IsSectionActive, _ := session.SessinStatus(c, "admin_session_cookie")
	if !IsSectionActive {
		fmt.Println("No Active Sessions found ")
		// c.HTML(http.StatusOK, "admin_login.html", []string{"a", "b", "c"})
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "Admin Login",
				"diplay": "none",
			},
		)
	} else {
		itemID := c.Request.URL.Query()["itemid"][0] // Getting Order ID passed with URL
		fmt.Println("Initiating to View/Edit item ,having ID", itemID)
		//populateCategoryItems(c, itemID)
		//GetItemDetails(itemID string) (itemDesc string, itemRate float64, unit string,itmID,itmStock int,itmBuyRate float64) {
		//Don't Confuse above function will redirect
		//to edit page, usual practice is giving here
		//but we achived this by modifying the existing
		//code so it happend so..
		itmDesc, itmSelRate, itmUnit, itmID, itmStock, itmBuyPrice := db.GetItemDetails(itemID)
		c.HTML(
			http.StatusOK,
			"edit_item.html", gin.H{
				"delWarning":   "none",
				"updateSucess": "none",
				"title":        "Edit Item",
				"itmID":        itmID,
				"itmDesc":      itmDesc,
				"itmUnit":      itmUnit,
				"itmBuyPrice":  itmBuyPrice,
				"itmSelRate":   itmSelRate,
				"itmStock":     itmStock,
			})
	}
}
func viewandedititemPost(c *gin.Context) {
	userOption := c.PostForm("action")
	itemID := c.PostForm("icode")
	switch userOption {
	case "delete":
		oK := db.DelitemFromMaster(itemID)
		if !oK { // Most of the case due to existing reference to other tables , this operation will fail
			itmDesc, itmSelRate, itmUnit, itmID, itmStock, itmBuyPrice := db.GetItemDetails(itemID)
			fmt.Println("DEL Operation Failed ")
			c.HTML(
				http.StatusOK,
				"edit_item.html", gin.H{
					"delWarning":   "block",
					"updateSucess": "none",
					"title":        "Edit Item",
					"itmID":        itmID,
					"itmDesc":      itmDesc,
					"itmUnit":      itmUnit,
					"itmBuyPrice":  itmBuyPrice,
					"itmSelRate":   itmSelRate,
					"itmStock":     itmStock,
				})
		} else { //Delete  success so going back to admin panel
			populateCategoryItems(c, "")
		}
	case "update":
		pPrice := c.PostForm("pPrice")
		sPrice := c.PostForm("sPrice")
		stock := c.PostForm("stock")
		itmName := c.PostForm("itmName")
		oK := db.UpdateItemInMaster(itemID, itmName, pPrice, sPrice, stock)
		if !oK {
			fmt.Println("Update Failed ")
		}
		//Diplaying after updation
		itmDesc, itmSelRate, itmUnit, itmID, itmStock, itmBuyPrice := db.GetItemDetails(itemID)
		c.HTML(
			http.StatusOK,
			"edit_item.html", gin.H{
				"delWarning":   "none",
				"updateSucess": "block",
				"title":        "Edit Item",
				"itmID":        itmID,
				"itmDesc":      itmDesc,
				"itmUnit":      itmUnit,
				"itmBuyPrice":  itmBuyPrice,
				"itmSelRate":   itmSelRate,
				"itmStock":     itmStock,
			})

	}
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
	router.GET("/viewandedititem", viewandedititemGet)
	router.POST("/viewandedititem", viewandedititemPost)
	//user Routes
	router.GET("/", userIndexGet)
	router.POST("/", userIndexPost)
	router.GET("/signup", userSignupGet)
	router.POST("/signup", userSignupPost)
	router.GET("/userlogin", userLoginGet)
	router.POST("/userlogin", userLoginPost)
	router.POST("/userotpverify", userOtpVerifyPost)
	router.GET("/orders", userOrdersGet)
	router.POST("/orders", userOrdersPost)
	router.GET("/usrlogout", userLogoutGet)
	router.POST("/orderconfirm", orderconfirmPost) //Once order Confirmed from /orders , it comes here
	router.POST("/orderhistory", orderHistoryPost)
	router.GET("/orderhistory", orderHistoryGet)    // to give functionality to Nav bar - Order Link
	router.GET("/showsingleorder", viewAnyOrderGet) // veiw any single Order

	// //TestCode
	// router.GET("/otp", otpGet)
	// router.GET("/sk",skGet)
	// router.GET("/gk",gkGet)
	router.Run()
}
