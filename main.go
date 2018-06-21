package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"flag"
)

type LoginResponse struct {
	Status  string
	Error   string
	Message string
	Data struct {
		AuthToken string
		UserId    string
	}
}

type PostMessageResponse struct {
	Success bool
	Error   string
}

func login(user, password, server string) LoginResponse {
	body := strings.NewReader(`user=` + user + `&password=` + password)
	req, err := http.NewRequest("POST", server+"/api/v1/login", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	loginResponse := LoginResponse{}
	err = json.Unmarshal(b, &loginResponse)
	return loginResponse
}

func postMessage(channel, message, userToken, userId, server string) PostMessageResponse {
	body := strings.NewReader(`channel=#` + channel + `&text=` + message)
	req, err := http.NewRequest("POST", server+"/api/v1/chat.postMessage", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Auth-Token", userToken)
	req.Header.Set("X-User-Id", userId)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	postMessageResponse := PostMessageResponse{}
	err = json.Unmarshal(b, &postMessageResponse)
	return postMessageResponse
}

func main() {
	channel := flag.String("c", "general", "Channel used to post the message")
	message := flag.String("m", "Hey, I'm a message!", "Message to post")
	user := flag.String("u", "user", "Rocket.Chat user")
	password := flag.String("p", "password", "Rocket.Chat user's password")
	server := flag.String("s", "http://localhost:3000", "Rocket.Chat server")
	flag.Parse()

	loginData := login(*user, *password, *server)
	if loginData.Status == "success" {
		postMessageData := postMessage(*channel, *message, loginData.Data.AuthToken, loginData.Data.UserId, *server)
		if postMessageData.Success {
			log.Println("Message sent to channel: #" + *channel + " <"+*server+"> message: " + *message + " [ok]")
		} else {
			log.Fatal(postMessageData.Error)
		}
	} else {
		log.Fatal(loginData.Error)
	}
}
