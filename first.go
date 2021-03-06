package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
				case "street_number":
					fallthrough
				case "route":
					message.Address += elemResults.Address_Components[key].Long_Name
				case "locality":
					message.City = elemResults.Address_Components[key].Long_Name
				case "administrative_area_level_1":
					message.Code = elemResults.Address_Components[key].Short_Name
					message.Name = elemResults.Address_Components[key].Long_Name
				case "postal_code":
					code, _ := strconv.Atoi(elemResults.Address_Components[key].Long_Name)
					message.Zipcode = code
				case "country":
					message.Country = elemResults.Address_Components[key].Long_Name
				}

			}
		}
	}
	return message
}

const login = "demo:demo1"

func main() {

	EnableLog := flag.Bool("log", false, "Enable log")
	flag.Parse()
	if !*EnableLog {
		log.SetOutput(ioutil.Discard)
	}

	StateName := map[string]string{
		"fl": "Florida",
	}

	http.HandleFunc("/address/normalize", func(w http.ResponseWriter, r *http.Request) {

		//---------------------------------------

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted Area"`)

		log.Println("WWW-Authenticate")

		slice := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		fmt.Println(len(slice))
		if len(slice) != 2 {
			log.Println("Not authorized, http error 401")
			http.Error(w, "Not authorized", 401)
			return
		}

		encodelogin, err := base64.StdEncoding.DecodeString(slice[1])
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), 401)
			return
		}

		if string(encodelogin) != login {
			log.Println("Bad login or passwrd")
			http.Error(w, "Not authorized", 401)
			return
		}

		log.Println("Authorized user")

		//---------------------------------------

		query := r.URL.Query()

		log.Print("Query: ")
		log.Println(query)

		log.Println(`Search "raw_address"`)

		raw_address, ok := query["raw_address"]

		if !ok {
			log.Println(`"raw_address" not found`)
			sendResponse(w, 404, MsgErr{"raw_address required"})
			return
		}

		log.Println(`"raw_address" found`)

		part := strings.Split(raw_address[0], ",")

		//--------------------------------------
		strAddress := strings.Replace(raw_address[0], " ", "+", -1)

		log.Println("Address search in Google")

		log.Println("http.Get: ", "https://maps.googleapis.com/maps/api/geocode/json?address=", strAddress, "&key=AIzaSyC-OyuXWSaNdtjcCTC4oz7W1jxv5MwCP8k&language=en")
		respGoogle, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + strAddress + "&key=AIzaSyC-OyuXWSaNdtjcCTC4oz7W1jxv5MwCP8k&language=en")

		var resultGoogle GoogleResponse
		err = json.NewDecoder(respGoogle.Body).Decode(&resultGoogle)
		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(w, "Error: %v", err)
		}

		log.Println("Address found")

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

		log.Println(`Found 5 part`)

		if len(part) == 5 {
			log.Println(`Сreate a message (json)`)
			code, _ := strconv.Atoi(part[3])
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			msg := Message{part[0], part[1], State{part[2], StateName[strings.Trim(part[2], " ")]}, code, part[4]}
			response, _ := json.Marshal(msg)
			fmt.Fprintf(w, "%+v", string(response))
		} else {
			log.Println(`The message (json) is not created`)
		}

		log.Println(`Тhe message (json) is created`)

	})
	log.Println("listen tcp localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
