package main

import (
	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
)

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
}

const (
	sdkClientName    = "my-amazing-app"
	sdkClientVersion = "2.0.5"
	credFilePath     = ".devenv.json"
)

func main() {
	cred := ReadJSONCredentials(credFilePath)
	cred.ClientName = sdkClientName
	cred.ClientVersion = sdkClientVersion

	log := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(log)      // Set the logger as default for TTN

	// Create a new SDK configuration for the public community network
	config := ttnsdk.NewCommunityConfig(cred.ClientName)
	config.ClientVersion = "2.0.5" // The version of the application

	client := config.NewClient(cred.TTNAppID, cred.TTNAppAccessKey)

	defer client.Close()

	// Start Publish/Subscribe client (MQTT)
	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not get application pub/sub", cred.ClientName)
	}

	defer pubsub.Close()

	// Get a publish/subscribe client scoped to my-test-device
	myNewDevicePubSub := pubsub.Device(cred.TTNDeviceID)

	defer myNewDevicePubSub.Close()

	// Subscribe to uplink messages
	uplink, err := myNewDevicePubSub.SubscribeUplink()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not subscribe to uplink messages", cred.ClientName)
	}

	for message := range uplink {
		log.Infof("%s: received uplink %s", cred.ClientName, JSONStringify(message.PayloadFields))
	}
}
