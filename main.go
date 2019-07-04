package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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

func loadConfigFromEnv(user, password, server, channel *string) {
	envUser := os.Getenv("ROCKET_CHAT_USER")
	if *user == "" && envUser != "" {
		*user = envUser
	}

	envPassword := os.Getenv("ROCKET_CHAT_PASSWORD")
	if *password == "" && os.Getenv("ROCKET_CHAT_PASSWORD") != "" {
		*password = envPassword
	}

	envServer := os.Getenv("ROCKET_CHAT_SERVER")
	if *server == "" && os.Getenv("ROCKET_CHAT_SERVER") != "" {
		*server = envServer
	}

	envChannel := os.Getenv("ROCKET_CHAT_CHANNEL")
	if *channel == "" && os.Getenv("ROCKET_CHAT_CHANNEL") != "" {
		*channel = envChannel
	}
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
		log.Println("Loading configuration from file " + *fromFile)
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

	loadConfigFromEnv(user, password, server, channel)

	if *user == "" || *password == "" || *server == "" || *channel == "" {
		log.Fatal("Please provide all the needed params to execute the application. Use rocket-notification -h to read the help.")
	}

	if *message == "" {
		log.Println("Reading text from stdin")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		newMsg := ""
		for err == nil {
			newMsg = newMsg + line
			line, err = reader.ReadString('\n')
		}
		message = &newMsg
		fmt.Println(*message)
	}

	if *message == "" {
		log.Fatal("no message was set")
	}

	if *isCode {
		newMsg := "```" + *message + "```"
		message = &newMsg
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
