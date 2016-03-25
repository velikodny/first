package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	//"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Message struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   `json:"state"`
	Zipcode int    `json:"zipcode"`
	Country string `json:"country"`
}

type GoogleResponse struct {
	Results []Results
	Status  string
}

type Results struct {
	Address_Components []Components
}

type Components struct {
	Long_Name  string
	Short_Name string
	Types      []string
}

type State struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type MsgErr struct {
	Err string `json:"error"`
}

func sendResponse(w http.ResponseWriter, status int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func bypassResultGoogle(res []Results) Message {
	message := Message{}
    
    for _, elemResults := range res {
		for key, elemAddress := range elemResults.Address_Components {
			for _, elemTypes := range elemAddress.Types {
				switch elemTypes {
                case "street_number": fallthrough
                case "route" : 
                        message.Address += elemResults.Address_Components[key].Long_Name 
                case "locality" :
                         message.City = elemResults.Address_Components[key].Long_Name 
                case "administrative_area_level_1" : 
                        message.Code = elemResults.Address_Components[key].Short_Name
                        message.Name = elemResults.Address_Components[key].Long_Name                        
                case "postal_code" : 
                        code, _ := strconv.Atoi(elemResults.Address_Components[key].Long_Name) 
                        message.Zipcode = code
                case "country" : 
                        message.Country = elemResults.Address_Components[key].Long_Name    
                }

			}
		}
	}
	return message
}

const login = "demo:demo1"

func main() {

	StateName := map[string]string{
		"fl": "Florida",
	}

	http.HandleFunc("/address/normalize", func(w http.ResponseWriter, r *http.Request) {

		//---------------------------------------

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted Area"`)

		slice := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		fmt.Println(len(slice))
		if len(slice) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		encodelogin, err := base64.StdEncoding.DecodeString(slice[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if string(encodelogin) != login {
			http.Error(w, "Not authorized", 401)
			return
		}

		//---------------------------------------

		query := r.URL.Query()

		raw_address, ok := query["raw_address"]

		if !ok {
			sendResponse(w, 404, MsgErr{"raw_address required"})
			return
		}

		part := strings.Split(raw_address[0], ",")

		//--------------------------------------
		strAddress := strings.Replace(raw_address[0], " ", "+", -1)

        respGoogle, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + strAddress + "&key=AIzaSyC-OyuXWSaNdtjcCTC4oz7W1jxv5MwCP8k&language=en")

		var resultGoogle GoogleResponse
		err = json.NewDecoder(respGoogle.Body).Decode(&resultGoogle)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}

		//msgPass := Message{}

	
			findField := bypassResultGoogle(resultGoogle.Results)

			fmt.Println(findField)

		//---------------------------------------

		for i, elem := range part {
			part[i] = strings.Trim(elem, " ")
		}

		if len(part) != 5 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(MsgErr{"raw_address required 2"})
			return
		}

		if len(part) == 5 {
			code, _ := strconv.Atoi(part[3])
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			//w.Header().Set("Content-Type", "application/json; charset=utf-8")
			msg := Message{part[0], part[1], State{part[2], StateName[strings.Trim(part[2], " ")]}, code, part[4]}
			response, _ := json.Marshal(msg)
			fmt.Fprintf(w, "%+v", string(response))
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
