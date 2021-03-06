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
var Dbug = true

//Connection .. to refer in all modules
var Connection *sql.DB

//Admin ... to pic admin details from db
type Admin struct {
	id                     int
	name, mail, sesid, pwd string
}

type order struct {
	ID       int
	Date     string
	CustName string
	Amt      float64
	Status   string
	PayMode  string
}

type ord struct { // using only while fetching single order details for user
	ID       int
	Desc     string
	Price    float64
	Unit     string
	Qty      int
	SubTotal float64 // This is must while listing single order details
}
type itemsInSingleOrder []ord

//Orders ... to collect Order details to show in order History of user.
type Orders []order

//OrderHistory ...
var OrderHistory Orders

//Connect for connecting with MariaDB
func Connect() {
	dbDriver := "mysql"
	dbUser := "fradmin"
	dbPass := "f#$@rider"
	dbName := "db_fishrider"
	var err error
	Connection, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println("dm-Error")
		panic(err.Error())
	} else if Dbug { //if in debug ON
		err = Connection.Ping() // Testing DB Connected or not
		if err != nil {
			fmt.Println("dm-db.Ping failed:", err)
		} else {
			fmt.Println("dm-db.Ping Succecss:")
		}
	}
}

//AddNewCustomer ... In final desing may return newly generated ID
func AddNewCustomer(mob, name, hName, sName, lMark string) error {
	insForm, err := Connection.Prepare(
		"INSERT INTO cust_master(cust_mob,cust_name,cust_adr1,cust_adr2,cust_lmark) VALUES (?,?,?,?,?)",
	)
	_, err = insForm.Exec(mob, name, hName, sName, lMark)
	if err != nil {
		fmt.Println(err)
		return err // If insert failed return error Only chance is while trying for duplicate entry

	}
	return nil

}

//InsertNewCatagory ...
func InsertNewCatagory(itemName, itemUnit string) (string, error) {
	insForm, err := Connection.Prepare(
		"INSERT INTO item_master(item_desc,item_unit) VALUES (?,?)",
	)
	_, err = insForm.Exec(itemName, itemUnit)
	if err != nil {
		return "", err // If insert failed return error Only chance is while trying for duplicate entry
	}

	//Geting last crated ID
	selDB, err := Connection.Query("SELECT max(item_id) FROM item_master")
	if err != nil {
		panic(err.Error())
	}
	var itemID int
	for selDB.Next() {
		err = selDB.Scan(&itemID)
		if err != nil {
			panic(err.Error())
		}
	}

	return strconv.Itoa(itemID), nil
}

//TraceAdminWithSIDinDB ... locate user with sessID
func TraceAdminWithSIDinDB(sessCookie string) (sessStatus bool, adminName string) {
	selDB, err := Connection.Query("SELECT admin_uname FROM admin_master WHERE admin_sesid = '" + sessCookie + "' ")
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&adminName)
		if err != nil {
			panic(err.Error())
		}
	}
	if adminName != "" {
		sessStatus = true
		return sessStatus, adminName
	}
	return sessStatus, adminName
}

//TraceTempSIDinDB ...
func TraceTempSIDinDB(sessCookie string) (sessStatus bool) {
	selDB, err := Connection.Query("SELECT temp_sesid FROM temp_cart WHERE temp_sesid = '" + sessCookie + "' ")
	var tempSesid string // Just storing and later will check whether it is existing or not
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&tempSesid)
		if err != nil {
			panic(err.Error())
		}
	}
	if tempSesid != "" {
		sessStatus = true
		return
	}
	return //by default sessStatus value will be false , so we are safe
}

//TraceUserWithSIDinDB ...
func TraceUserWithSIDinDB(sessCookie string) (sessStatus bool, custMob int, custName, custAdr1, custAdr2 string, custID int) {
	selDB, err := Connection.Query("SELECT cust_name,cust_mob,cust_adr1,cust_adr2,cust_id FROM cust_master WHERE cust_sesid = '" + sessCookie + "' ")
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&custName, &custMob, &custAdr1, &custAdr2, &custID)
		if err != nil {
			panic(err.Error())
		}
	}
	if custName != "" {
		sessStatus = true
		return
	}
	return
}

//GetLastItemID ...  Item ID is the name of slecet tag in html, so
//at preset not found any other logic to pic the selected items from
//item listing page , TBD : find out some other logic later
func GetLastItemID() (count int) {
	count = 0
	selDB, err := Connection.Query("SELECT MAX(item_id) FROM item_master")
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&count)
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

//GetItemDetails ... to prepate structure for rendering
func GetItemDetails(itemID string) (itemDesc string, itemRate float64, unit string, itmID, itmStock int, itmBuyRate float64) {
	//item_id,,,item_stock,,item_buy_price
	selDB, err := Connection.Query("SELECT item_desc,item_sel_price,item_unit,item_id,item_stock,item_buy_price FROM item_master WHERE item_id=" + itemID)
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&itemDesc, &itemRate, &unit, &itmID, &itmStock, &itmBuyRate)
		if err != nil {
			panic(err.Error())
		}
	}
	return
}

//AddNewOrderEntry ...
func AddNewOrderEntry(custID int, orderAmt float64, modeOfPay string) (operationStatus bool, generatedOrderID int) {
	insForm, err := Connection.Prepare(
		"INSERT INTO order_master(cust_id,order_amt,p_mode) VALUES (?,?,?)",
	)
	_, err = insForm.Exec(custID, orderAmt, modeOfPay)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Geting last crated ID
	selDB, err := Connection.Query("SELECT max(order_id) FROM order_master")
	if err != nil {
		fmt.Println(err)
		return
	}
	for selDB.Next() {
		err = selDB.Scan(&generatedOrderID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

//AddNewOrderDetails ...
func AddNewOrderDetails(orderID, itemID, itemQty int) (operationStatus bool) {
	insForm, err := Connection.Prepare(
		"INSERT INTO order_detail(order_id,item_id,item_qty) VALUES (?,?,?)",
	)
	_, err = insForm.Exec(orderID, itemID, itemQty)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Hey we need to reduce the stock too
	insForm, err = Connection.Prepare(
		"UPDATE item_master SET item_stock = item_stock - ? WHERE item_id = ?",
	)
	_, err = insForm.Exec(itemQty, itemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Everything looks fine!!
	operationStatus = true
	return
}

//GetOrderHistory ...
func GetOrderHistory(custID int) (operationStatus bool, UsrOrderHistory Orders) {
	var singleOrder order
	var strSQL string
	if custID > 0 { //pic details of this user
		fmt.Println("Tarcing Oder hostory of Custer having ID:", custID)
		strSQL = "SELECT order_id,order_date,order_amt,order_status,p_mode from order_master WHERE cust_id= " + strconv.Itoa(custID)
	}else if(custID == 0){ // Pick all to show in Order view of Admin
		strSQL = "SELECT order_master.order_id,order_master.order_date,cust_master.cust_name,order_master.order_amt,order_master.order_status,order_master.p_mode FROM order_master INNER JOIN cust_master ON order_master.cust_id=cust_master.cust_id WHERE order_master.order_status='pending'"
	}else { // pic all pending orders
		strSQL = "SELECT order_master.order_id,order_master.order_date,cust_master.cust_name,order_master.order_amt,order_master.order_status,order_master.p_mode FROM order_master INNER JOIN cust_master ON order_master.cust_id=cust_master.cust_id"
	}

	selDB, err := Connection.Query(strSQL)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Here also based on rquest ( need all pending OR only one based on CustID) , Parameters in Scan also differe
	if custID > 0 {
		for selDB.Next() {
			err = selDB.Scan(&singleOrder.ID, &singleOrder.Date, &singleOrder.Amt, &singleOrder.Status, &singleOrder.PayMode)
			if err != nil {
				fmt.Println(err)
				return
			}
			UsrOrderHistory = append(UsrOrderHistory, singleOrder)
		}
	} else {
		for selDB.Next() {
			err = selDB.Scan(&singleOrder.ID, &singleOrder.Date, &singleOrder.CustName, &singleOrder.Amt, &singleOrder.Status, &singleOrder.PayMode)
			if err != nil {
				fmt.Println(err)
				return
			}
			UsrOrderHistory = append(UsrOrderHistory, singleOrder)
		}
	}
	operationStatus = true // all went well ..
	return
}

//GetSingleOredrDetails ...
func GetSingleOredrDetails(OrderID string) (operationStatus bool, ItemsInOrder itemsInSingleOrder, date, ordStatus, PayMode string, amt float64) {
	selDB, err := Connection.Query("SELECT order_date,order_status,order_amt,p_mode from order_master WHERE order_id= " + OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for selDB.Next() {
		err = selDB.Scan(&date, &ordStatus, &amt, &PayMode)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//Now Collect Multtiple Items ( maye contain single item only in one order :) , in this order
	var item ord // we are picking one by one. adding to the struct array
	selDB, err = Connection.Query("select order_detail.item_id,item_master.item_desc,item_master.item_sel_price,item_master.item_unit,order_detail.item_qty FROM order_master INNER JOIN order_detail ON order_master.order_id = order_detail.order_id INNER JOIN item_master ON order_detail.item_id = item_master.item_id WHERE order_master.order_id=" + OrderID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for selDB.Next() {
		err = selDB.Scan(&item.ID,
			&item.Desc,
			&item.Price,
			&item.Unit,
			&item.Qty,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		item.SubTotal = (item.Price * float64(item.Qty))
		ItemsInOrder = append(ItemsInOrder, item) // making a collection of all items in our oder
	}
	// okay done the job
	operationStatus = true
	return
}

//DelitemFromMaster ...
func DelitemFromMaster(itemID string) (oprationStatus bool) {
	delForm, err := Connection.Prepare("DELETE FROM item_master WHERE item_id= ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = delForm.Exec(itemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	oprationStatus = true
	return
}

//UpdateItemInMaster ...
func UpdateItemInMaster(itemID, itemDesc string, Pprice, Sprice string, stock string) (oprationStatus bool) {
	updateForm, err := Connection.Prepare("UPDATE item_master SET item_desc=?,item_buy_price=?,item_sel_price=?,item_stock=? WHERE item_id=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = updateForm.Exec(itemDesc, Pprice, Sprice, stock, itemID)
	if err != nil {
		fmt.Println(err)
		return
	}
	oprationStatus = true
	return
}

//UpdateOrderStatusToApproved ...
func UpdateOrderStatusToApproved(ordID string) (operationStatus bool) {
	updateForm, err := Connection.Prepare("UPDATE order_master SET order_status='approved' WHERE order_id=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = updateForm.Exec(ordID)
	if err != nil {
		fmt.Println(err)
		return
	}
	operationStatus = true
	return

}
