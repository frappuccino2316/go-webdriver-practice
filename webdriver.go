package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type WebDriver struct {
	SessionID string
}

type sessionResponse struct {
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

	startUrl := "http://localhost:9515/session"

	cap := `
{
	"capabilities": {
		"alwaysMatch": {
			"goog:chromeOptions": {
				"args": ["--no-sandbox"]
			}
		}
	}
}
`

	req, err := http.NewRequest("POST", startUrl, strings.NewReader(cap))

	if err != nil {
		return WebDriver{}, err
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return WebDriver{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ReadAll Error! %v\n", err)
		return WebDriver{}, err
	}

	var data sessionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return WebDriver{}, err
	}

	time.Sleep(time.Second * 5)

	driver := WebDriver{SessionID: data.Value.SessionID}
	return driver, nil
}
