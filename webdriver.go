package main

import (
	"net/http"
	"os/exec"
	"runtime"

	"prawd/api"
)

type WebDriver struct {
	SessionID string
	Client    *http.Client
}

type SessionResponse struct {
	Value struct {
		SessionID    string      `json:"sessionId"`
		Capabilities interface{} `json:"capabilities"`
	} `json:"value"`
}

func NewDriver() (WebDriver, error) {
	var binaryName string
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		binaryName = "chromedriver.exe"
		command = exec.Command("start", binaryName)
	} else {
		binaryName = "chromedriver"
		command = exec.Command(binaryName)
	}

	if err := command.Start(); err != nil {
		return WebDriver{}, err
	}

	client := new(http.Client)

	sessionID, err := api.NewSession(client)

	if err != nil {
		return WebDriver{}, nil
	}

	return WebDriver{SessionID: sessionID, Client: client}, nil
}

func (wd WebDriver) Close() error {
	err := api.DeleteSession(wd.Client, wd.SessionID)

	if err != nil {
		return err
	}

	return nil
}
