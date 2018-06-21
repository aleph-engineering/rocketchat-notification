package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

const SERVER = "http://192.168.1.43:3030"
const USER = "jenkins"
const PASSWORD = "owulacja"
const CHANNEL = "jenkins"
const MESSAGE = "This is a friendly test messaje"

func TestLogin(t *testing.T) {
	loginData := login(USER, PASSWORD, SERVER)
	assert.Equal(t, "success", loginData.Status)
	assert.NotNil(t, loginData.Data)
}

func TestLoginFailed(t *testing.T) {
	loginData := login(USER+"false", PASSWORD, SERVER)
	assert.Equal(t, "error", loginData.Status)
}

func TestPostMessage(t *testing.T) {
	loginData := login(USER, PASSWORD, SERVER)
	assert.Equal(t, "success", loginData.Status)
	if loginData.Status == "success" {
		postMessageData := postMessage(CHANNEL, MESSAGE, loginData.Data.AuthToken, loginData.Data.UserId, SERVER)
		assert.True(t, postMessageData.Success)
	} else {
		log.Fatal(loginData.Error)
	}
}
