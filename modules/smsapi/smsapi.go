package smsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-arun/fishrider/modules/db"
)

type sendResponse struct { // to store the values receving while initating for OTP, and Otp_id is required later in varification process
	Otp_id, Status string
	expiry         int
}
type verifyResponse struct { // to store response of OTP verification with code received in Mobile
	Status, Error string
}

//GenerateOTP ...
func GenerateOTP(mob string) (optID string) {
	return "comment this line later"
	mob = "91"+mob // added countryCode
	fmt.Println("Sening OTP to Mobile:",mob)
	var apiInitalResponse sendResponse

	url := "https://d7networks.com/api/verifier/send"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("mobile", mob)
	_ = writer.WriteField("sender_id", "SMSINFO")
	_ = writer.WriteField("message", "Your otp code is {code}")
	_ = writer.WriteField("expiry", "900")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Token eab00d735d6a78c77539bcb0389e36fc62bde2e6")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	json.Unmarshal(body, &apiInitalResponse)
	fmt.Println("OTP-ID-->", apiInitalResponse.Otp_id)
	return apiInitalResponse.Otp_id

}

//VerifyOTP ...
func VerifyOTP(otpID, otpCode string) (status bool, respMsg string) {
	respMsg = "Comment below 3 liens later"
	status = true
	return

	var apiResponse verifyResponse

	url := "https://d7networks.com/api/verifier/verify"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("otp_id", otpID)
	_ = writer.WriteField("otp_code", otpCode)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Token e502106fa304c9d842f357ceee0bc500b6393109")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &apiResponse)
	respMsg = apiResponse.Error
	fmt.Println(apiResponse.Status, apiResponse.Error)
	if apiResponse.Status == "success" {
		status = true
	}
	if db.Dbug {
		fmt.Println("OTP Verification Status :", string(body))
	}
	return
}
