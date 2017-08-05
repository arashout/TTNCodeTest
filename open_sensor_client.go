package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// OpenSensorClient ...
type OpenSensorClient struct {
	httpClient  http.Client
	credentials *Credentials
}

// InitializeOpenSensorClient ...
// Pretty simple function that returns a OpenSensorClient pointer
func InitializeOpenSensorClient(cred *Credentials) *OpenSensorClient {
	osc := &OpenSensorClient{
		httpClient:  http.Client{},
		credentials: cred,
	}
	return osc
}

// SendDataToTopic ...
// Function to send provided json data to the topic url specified
// in the credentials
func (osc *OpenSensorClient) SendDataToTopic(jsonData string) {
	valueWrapper := map[string]string{"data": jsonData}
	jsonPayload, err := json.Marshal(valueWrapper)
	if err != nil {
		log.Fatalf("Unable to marshal jsonData: %s", err.Error())
	}
	requestURL := osc.buildRequestURL()

	// Build Request
	req, err := http.NewRequest("POST", requestURL.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Failed to create request: %s", err.Error())

	}
	req.Header.Set("Authorization", fmt.Sprintf("api-key %s", osc.credentials.OpenSensorAPIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := osc.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send POST request: %s", err.Error())
		return
	}

	log.Printf("POST request of %s to %s returned with status: %s", jsonData, requestURL.Path, resp.Status)
}

func (osc *OpenSensorClient) buildRequestURL() *url.URL {
	u, err := url.Parse(osc.credentials.OpenSensorTopicURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("client-id", osc.credentials.OpenSensorClientID)
	q.Set("password", osc.credentials.OpenSensorClientPassword)
	u.RawQuery = q.Encode()
	return u
}
