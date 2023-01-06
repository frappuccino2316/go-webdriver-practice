package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type sessionResponse struct {
	Value struct {
		SessionID    string      `json:"sessionId"`
		Capabilities interface{} `json:"capabilities"`
	} `json:"value"`
}

func NewSession(client *http.Client) (string, error) {
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
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ReadAll Error! %v\n", err)
		return "", err
	}

	var data sessionResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.Value.SessionID, nil
}

func DeleteSession(client *http.Client, sessionID string) error {
	closeURL := fmt.Sprintf("http://localhost:9515/session/%v", sessionID)

	params := "{}"
	req, err := http.NewRequest("POST", closeURL, strings.NewReader(params))

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
