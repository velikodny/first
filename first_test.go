package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBasicAuth(t *testing.T) {

	var user string = "home"
	var passwd string = "home"
	client := &http.Client{}
	req, err := http.NewRequest("GET", ":8080/address/normalize?raw_address=8150+nw+53rd+street+suite+350-140,+doral,+fl,+33166,usa", nil)
	req.SetBasicAuth(user, passwd)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil || resp.Body == nil{
		t.Error("Err ", err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Println(err)

	t.Error("Body ", s)

}