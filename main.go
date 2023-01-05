package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	var binaryName string
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		binaryName = "chromedriver.exe"
		command = exec.Command("start", binaryName)
	} else {
		binaryName = "chromedriver"
		command = exec.Command(binaryName)
	}

	fmt.Println("run & start")
	if err := command.Start(); err != nil {
		fmt.Println("error start")
		return
	}
	defer kill(command)

	fmt.Println("wait start")

	time.Sleep(time.Second * 5)

	// ===post===

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
		fmt.Printf("Request Error! %v\n", err)
		return
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Post Error! %v\n", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ReadAll Error! %v\n", err)
		return
	}

	var data SessionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Unmarshal Error! %v\n", err)
		return
	}
	fmt.Println(data)
	fmt.Println(data.Values.SessionId)

	time.Sleep(time.Second * 5)

	// 	// =========close=========

	closeUrl := fmt.Sprintf("http://localhost:9515/session/%v", data.Values.SessionId)

	closeRequestBody := map[string]string{}

	closeJsonString, err := json.Marshal(closeRequestBody)
	if err != nil {
		fmt.Printf("Encoding Error! %v\n", err)
	}

	req, err = http.NewRequest("DELETE", closeUrl, bytes.NewBuffer(closeJsonString))

	if err != nil {
		fmt.Printf("Request Error! %v\n", err)
		return
	}

	res, err = client.Do(req)
	if err != nil {
		fmt.Printf("Post Error! %v\n", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("finish")
}

func kill(command *exec.Cmd) {
	if runtime.GOOS == "windows" {
		command.Process.Kill()
	} else {
		command.Process.Signal(syscall.SIGTERM)
	}
}

type SessionResponse struct {
	Values Value `json:"value"`
}

type Value struct {
	// Capabilities
	SessionId string `json:"sessionId"`
}

// type Capabilities struct {
// 	acceptInsecureCerts       bool
// 	browserName               string
// 	browserVersion            string
// 	chrome                    map[string]string
// 	goog                      map[string]string `json:"goog:chromeOptions"`
// 	networkConnectionEnabled  bool
// 	pageLoadStrategy          string
// 	platformName              string
// 	proxy                     map[string]string
// 	setWindowRect             bool
// 	strictFileInteractability bool
// 	timeouts                  map[string]int
// 	unhandledPromptBehavior   string
// 	credBlob                  bool `json:"webauthn:extension:credBlob"`
// 	largeBlob                 bool `json:"webauthn:extension:largeBlob"`
// 	virtualAuthenticators     bool `json:"webauthn:virtualAuthenticators"`
// }
