// Package session ...  is to handle session-id and passwd-hashing
package session

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-arun/fishrider/modules/db"
	"golang.org/x/crypto/bcrypt"
)

type item struct {
	ICode, IQty string
}
//CartItems ... exported because it is using in main too
type CartItems []item

var cart CartItems // All items selectig by user will be appened to this array and  later push it as cookie

//GenerateNewSessionID ... generating New UUID
func GenerateNewSessionID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

//AddAdminSessionIDToDB ... to coresponnding users's field in db
func AddAdminSessionIDToDB(userName string) string {
	sessID, _ := GenerateNewSessionID()

	insForm, err := db.Connection.Prepare("UPDATE admin_master SET admin_sesid=? WHERE admin_uname= ?")
	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(sessID, userName)
	return sessID
}

//AddUserSessionIDToDB ...
func AddUserSessionIDToDB(mobNo int) string {
	sessID, _ := GenerateNewSessionID()

	insForm, err := db.Connection.Prepare("UPDATE cust_master SET cust_sesid=? WHERE cust_mob= ?")
	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(sessID, mobNo)
	return sessID
}

//RemoveAdminSessionIDFromDB ...
func RemoveAdminSessionIDFromDB(c *gin.Context) {
	if db.Dbug {
		fmt.Println("Removing Session Id from DB and Browsser ...")
	}
	sessionCookie, _ := c.Cookie("admin_session_cookie")

	insForm, err := db.Connection.Prepare("UPDATE admin_master SET admin_sesid='' WHERE admin_sesid = ?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(sessionCookie)
	// Deleting cookie
	c.SetCookie("admin_session_cookie",
		"",
		-1, // delete now !!
		"/",
		"", false, false,
	)
}

//RemoveUserSessionIDFromDB ...
func RemoveUserSessionIDFromDB(c *gin.Context) {
	if db.Dbug {
		fmt.Println("Removing Session Id from DB and Browsser ...")
	}
	sessionCookie, _ := c.Cookie("user_session_cookie")

	insForm, err := db.Connection.Prepare("UPDATE cust_master SET cust_sesid='' WHERE cust_sesid = ?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(sessionCookie)
	// Deleting cookie
	c.SetCookie("user_session_cookie",
		"",
		-1, // delete now !!
		"/",
		"", false, false,
	)
}

//HashPassword ..
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//AdminCredentialsVerify ...
func AdminCredentialsVerify(uname, pwd string) bool {
	selDB, err := db.Connection.Query("SELECT admin_passwd,admin_sesid FROM admin_master WHERE admin_uname = '" + uname + "' ")
	if err != nil {
		panic(err.Error())
	}
	var sesID string
	var pwdInDB string // PasswdFromDB
	for selDB.Next() {
		err = selDB.Scan(&pwdInDB, &sesID)
		if err != nil {
			panic(err.Error())
		}
	}
	//Need to hashcheck with supplied passwd by user and the one in DB
	if pwdInDB != "" {
		result := CheckPasswordHash(pwd, pwdInDB)
		fmt.Println("PasswdMached-->", result)
		return result
	} else {
		return false // No Such user
	}
}

//UserCredentialsVerify ...
func UserCredentialsVerify(mobNum string) (userExist bool, custID int, custName string) {
	selDB, err := db.Connection.Query("SELECT cust_id,cust_name FROM cust_master WHERE cust_mob= '" + mobNum + "'")
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		err = selDB.Scan(&custID, &custName)
		if err != nil {
			panic(err.Error())
		}
	}
	//If any matching record found
	if custName == "" {
		return
	} else {
		userExist = true
		return
	}
}

//SessinStatus ...return true if there is an active admin session
func SessinStatus(c *gin.Context, cookieName string) (sesStatus bool) {
	sessionCookie, _ := c.Cookie(cookieName)
	if sessionCookie == "" { // no cookie found
		return sesStatus // by default 'recordFound' val will be false
	} //cookie received frim browser still we need to ensure from database too
	if cookieName == "admin_session_cookie" { // divert here based on which sesid user/admin
		sesStatus, _ = db.TraceAdminWithSIDinDB(sessionCookie)
	} else if cookieName == "user_session_cookie" {
		sesStatus, _,_,_,_ = db.TraceUserWithSIDinDB(sessionCookie)
	} else { // then it is about temp_sesid go ahed
		sesStatus = db.TraceTempSIDinDB(sessionCookie)
	}

	fmt.Println(cookieName, " Session Exists Status --> ?", sesStatus)
	return sesStatus //
}

//SetAdminSessionCookie ...
func SetAdminSessionCookie(c *gin.Context, adminUserName, sessionCookieName string) {

	sessionID := AddAdminSessionIDToDB(adminUserName) // Generate a new SSID and insert to DB
	c.SetCookie(sessionCookieName,
		sessionID,
		3600*12, // 12hrs
		"/",
		"", false, false, //domain excluded
	)
}

//SetUserSessionCookie ...
func SetUserSessionCookie(c *gin.Context, mobNumber int, sessionCookieName string) {

	sessionID := AddUserSessionIDToDB(mobNumber) // Generate a new SSID and insert to DB
	c.SetCookie(sessionCookieName,
		sessionID,
		3600*8760, // 8760hrs = 1 Year
		"/",
		"", false, false, //domain excluded
	)
}

//PushSelectionToCookie ...
func PushSelectionToCookie(c *gin.Context, itemCode, itemQty string) {
	var newItem item
	newItem.ICode = itemCode
	newItem.IQty = itemQty
	fmt.Println("newItem.iCode & newItem.iQty", newItem.ICode, newItem.IQty)
	cart = append(cart, newItem)
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	fmt.Println("Cartitem->", string(cartJSON))
	c.SetCookie("user_cart",
		string(cartJSON),
		3600*8760, // 8760hrs = 1 Year
		"/",
		"", false, false, //domain excluded
	)

}
//PullCartItemFromCookie ... Collect cart items from Cookie
func PullCartItemFromCookie(c *gin.Context )(CartItems){
	sessionCookie, _ := c.Cookie("user_cart")
	fmt.Println("sessionCookie-->",sessionCookie)
	var itemInCart CartItems
	json.Unmarshal([]byte(sessionCookie), &itemInCart)
	return itemInCart

	// for key,_ := range itemInCart{ // DONT REMOVE MAY REQ. LATER
	// 	fmt.Println(itemInCart[key].ICode)
	// 	// fmt.Println("key and val=",key,val)
	// }

}
