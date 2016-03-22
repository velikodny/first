package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"io"
)

type Message struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   `json:"state"`
	Zipcode int    `json:"zipcode"`
	Country string `json:"country"`
}

type MessageGoogle struct {
	Address string `json:"results.formatted_address"`
	City    string 
	State   
	Zipcode int    
	Country string 
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
        
		//---------------------------------------

        //sep := []byte(" ")
       // partGoogle := strings.Split(raw_address[0], " ")
       // addressForGoogle = strings.Join(partGoogle,"+")
       
        strAddress := strings.Replace(raw_address[0], " ", "+", -1)
		respGoogle, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + strAddress +"&key=AIzaSyC-OyuXWSaNdtjcCTC4oz7W1jxv5MwCP8k&language=en")

        var mesgGoogle MessageGoogle
       err = json.NewDecoder(respGoogle.Body).Decode(&mesgGoogle)
       err = json.NewDecoder(io.LimitReader(respGoogle.Body, 64)).Decode(&mesgGoogle)
       // err = json.Unmarshal(respGoogle, &mesgGoogle)   
        fmt.Println(mesgGoogle)
		
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
