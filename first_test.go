package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func Auth (user, passwd string) (*http.Response, error){
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/address/normalize?raw_address=8150+nw+53rd+street+suite+350-140,+doral,+fl,+33166,usa", nil)

	if (user != " " && passwd != " "){
		req.SetBasicAuth(user, passwd)	
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()	
	}
	
	return resp	, err
}

func CheckAuth(res *http.Response, err error) (string){
	if err != nil || res.Body == nil{
		return "Err " + res.Status	 
	}	
	
	bodyText, err := ioutil.ReadAll(res.Body)

	if res.StatusCode == http.StatusOK{
		return ""
	}

	if res.StatusCode != http.StatusUnauthorized{
		return "Err " + string(bodyText) +  "  " + res.Status
	}	
	
	return ""
}

func TestNotAuth(t *testing.T) {
	
	res, err := Auth(" ", " ")

	if err != nil || res == nil {
		t.Error(err)		
	}	

	if res.StatusCode != http.StatusUnauthorized{
		t.Error("Status ", res.Status)
	} 
	
}

func TestBasicAuth(t *testing.T) {

	res, err := Auth("demo", "demo1")
	if err != nil{
		t.Error(err)		
	}	
	if res.StatusCode != http.StatusOK{
		t.Error("Status ", res.Status)
	}	
}

func TestBasicAuthBlanck(t *testing.T) {

	res, err := Auth("", "")
	if err != nil{
		t.Error(err)		
	}	
	if res.StatusCode != http.StatusUnauthorized{
		t.Error("Status ", res.Status)
	}	
}

func TestMassage(t *testing.T) {

	res, err := Auth("demo", "demo1")
	if err != nil{
		t.Error(err)		
	}	
	if res.StatusCode != http.StatusOK{
		t.Error("Status ", res.Status)
	}	
	
	fmt.Printf("[%v]", res.Body == nil)
}