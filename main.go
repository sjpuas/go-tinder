package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var token_facebook = ""
var urlBase = ""

func main() {

	conf := readConf()
	token_facebook = conf.Token_facebook
	urlBase = conf.UrlBase

	authTinder := auth()
	fmt.Println("Hello ", authTinder.User.Full_name)
	recs := recs(authTinder)
	for _, rec := range recs.Results {
		match := like(rec._id, authTinder)
		fmt.Println("Name: ", rec.Name, " ,Match:", match)
	}

	fmt.Println("Finish")

}

func auth() AuthTinder {
	var authTinder AuthTinder
	var jsonStr = []byte(`{"facebook_token":"` + token_facebook + `"}`)
	req, err := http.NewRequest("POST", urlBase+"/auth", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &authTinder)
	return authTinder
}

func recs(authTinder AuthTinder) RecsTinder {
	var recs RecsTinder
	req, err := http.NewRequest("GET", urlBase+"/user/recs", nil)
	req.Header.Set("X-Auth-Token", authTinder.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &recs)
	return recs

}

func like(_id string, authTinder AuthTinder) bool {
	var like LikeTinder
	req, err := http.NewRequest("GET", urlBase+"/like/"+_id, nil)
	req.Header.Set("X-Auth-Token", authTinder.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &like)
	return like.Match

}

func readConf() Configuration {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}

type Configuration struct {
	Token_facebook string
	UrlBase        string
}

type UserTinder struct {
	_id        string
	Full_name  string
	Name       string
	Gender     int
	Bio        string
	Birth_date time.Time
	Photos     []PhotoTinder
}

type AuthTinder struct {
	Token string
	User  UserTinder
}

type RecsTinder struct {
	Status  int
	Results []UserTinder
}

type LikeTinder struct {
	Match bool
}

type PhotoTinder struct {
	Url            string
	ProcessedFiles []ProcessedFileTinder
}

type ProcessedFileTinder struct {
	Url    string
	Height int
	Width  int
}
