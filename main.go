package main

import "log"

// Credentials ... Type to store credentials
type Credentials struct {
	ClientName               string
	ClientVersion            string
	TTNAppID                 string `json:"ttnAppId"`
	TTNAppAccessKey          string `json:"ttnAppAccessKey"`
	TTNDeviceID              string `json:"ttnDeviceId"`
	OpenSensorAPIKey         string `json:"openSensorAPIKey"`
	OpenSensorClientID       string `json:"openSensorClientId"`
	OpenSensorClientPassword string `json:"openSensorClientPassword"`
	OpenSensorTopicURL       string `json:"openSensorTopicUrl"`
}

const (
	sdkClientName    = "TTN_Code_Test_App"
	sdkClientVersion = "2.0.5"
	credFilePath     = ".devenv.json"
	maxRequests      = 5
)

func main() {
	cred := ReadJSONCredentials(credFilePath)
	cred.ClientName = sdkClientName
	cred.ClientVersion = sdkClientVersion

	ttnClient := InitializeTTNClient(&cred)
	uplink := ttnClient.GetUplinkChannel()

	osc := InitializeOpenSensorClient(&cred)

	requestCount := 0
	for message := range uplink {
		jsonData := JSONStringify(message.PayloadFields)
		log.Printf("%s: received uplink %s", cred.ClientName, jsonData)
		go osc.SendDataToTopic(jsonData)

		// This part can be removed if you have a safe exit procedure
		// Maybe capture ctrl-c?
		requestCount = requestCount + 1
		if requestCount == maxRequests {
			break
		}
	}

	ttnClient.Close()
}
