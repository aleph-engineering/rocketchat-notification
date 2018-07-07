package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type LoginResponse struct {
	Status  string
	Error   string
	Message string
	Data    struct {
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
	message := flag.String("m", "", "Message to post")
	isCode := flag.Bool("code", false, "Wrap message in a code area (default false)")
	user := flag.String("u", "", "Rocket.Chat user")
	password := flag.String("p", "", "Rocket.Chat user's password")
	server := flag.String("s", "http://localhost:3000", "Rocket.Chat server")
	fromFile := flag.String("f", "", "Configuration file (Optional)")
	flag.Parse()

	if *fromFile != "" {
		config := ReadConfig(*fromFile)
		if config.User != "" && *user == "" {
			user = &config.User
		}
		if config.Password != "" && *password == "" {
			password = &config.Password
		}
		if config.Server != "" && *server == "http://localhost:3000" {
			server = &config.Server
		}
	}

	if *message == "" {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		new_msg := ""
		for err == nil {
			new_msg = new_msg + line
			line, err = reader.ReadString('\n')
		}
		message = &new_msg
	}

	if *isCode {
		new_msg := "```" + *message + "```"
		message = &new_msg
	}

	if *user == "" || *password == "" || *server == "" || *message == "" || *channel == "" {
		log.Fatal("Please provide all the needed params to execute the application. Use rocket-notification -h to read the help.")
	}

	loginData := login(*user, *password, *server)
	if loginData.Status == "success" {
		postMessageData := postMessage(*channel, *message, loginData.Data.AuthToken, loginData.Data.UserId, *server)
		if postMessageData.Success {
			log.Println("Message sent to channel: #" + *channel + " <" + *server + "> message: " + *message + " [ok]")
		} else {
			log.Fatal(postMessageData.Error)
		}
	} else {
		log.Fatal(loginData.Error)
	}
}
