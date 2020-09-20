// Package session ...  is to handle session-id and passwd-hashing
package session

import (
	"crypto/rand"
	"io"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-arun/fishrider/modules/db"
	"github.com/gin-gonic/gin"


)

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

//AddSessionIDToDB ... to coresponnding users's field in db
func AddSessionIDToDB(userName string)(string) {
	sessID,_ := GenerateNewSessionID()

	insForm, err := db.Connection.Prepare("UPDATE admin_master SET admin_sesid=? WHERE admin_name= ?")
	if err != nil {
		   panic(err.Error())
	}

insForm.Exec(sessID,userName)
    return sessID
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

//UsrCredentialsVerify ...
func UsrCredentialsVerify(uname,pwd string)(bool){
	selDB, err := db.Connection.Query("SELECT admin_passwd,admin_sesid FROM admin_master WHERE admin_name = '" + uname + "' ")
	if err !=nil {
			panic(err.Error())
	}
	var sesID string
	var pwdInDB string // PasswdFromDB
	for selDB.Next(){
		err = selDB.Scan(&pwdInDB, &sesID)
		if err != nil {
			panic(err.Error())
		}
	}
	//Need to hashcheck with supplied passwd by user and the one in DB
	if pwdInDB != "" {
		result := CheckPasswordHash(pwd,pwdInDB)
		fmt.Println("PasswdMached-->",result)
		return result
	}else{
		return false // No Such user
	}
	
}

//SessinStatus ...return true if there is an active admin session
func SessinStatus(c *gin.Context,cookieName string)(sesStatus bool) { 
	sessionCookie,_ := c.Cookie (cookieName)
	if sessionCookie == "" { // no cookie found 
		return sesStatus // by default 'recordFound' val will be false
	} //cookie received frim browser still we need to ensure from database too
	sesStatus,_ = db.TraceUserWithSID(sessionCookie)
	fmt.Println("Admin Session Exissts ?",sesStatus)
	return sesStatus // 
}

//SetSessionCookie ...
func SetSessionCookie(c *gin.Context,adminUserName,sessionCookieName string){

	sessionID := AddSessionIDToDB(adminUserName) // Generate a new SSID and insert to DB
	c.SetCookie(sessionCookieName,
	sessionID,
	3600*12, // 12hrs
	"/",
	"",false,false, //domain excluded 
	)
}


