package main

import (
	"log"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	"github.com/TheThingsNetwork/ttn/core/types"
)

// TTNClient ...
// TTNClient wrapper struct that allows me to bind methods easier
type TTNClient struct {
	client       ttnsdk.Client
	pubsub       ttnsdk.ApplicationPubSub
	devicePubSub ttnsdk.DevicePubSub
	credentials  *Credentials
}

// InitializeTTNClient ...
// Function that returns a TTNClient for client to use for method calls
// Example from https://godoc.org/github.com/TheThingsNetwork/go-app-sdk
func InitializeTTNClient(cred *Credentials) *TTNClient {
	logTNN := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(logTNN)      // Set the logger as default for TTN

	// Create a new SDK configuration for the public community network
	config := ttnsdk.NewCommunityConfig(cred.ClientName)
	config.ClientVersion = "2.0.5" // The version of the application

	var ttnClient TTNClient
	ttnClient.client = config.NewClient(cred.TTNAppID, cred.TTNAppAccessKey)

	// Start Publish/Subscribe client (MQTT)
	var err error
	ttnClient.pubsub, err = ttnClient.client.PubSub()
	if err != nil {
		logTNN.WithError(err).Fatalf("%s: could not get application pub/sub", cred.ClientName)
	}

	ttnClient.devicePubSub = ttnClient.pubsub.Device(cred.TTNDeviceID)

	return &ttnClient
}

// GetUplinkChannel ...
// This method return an Uplink channel of the device in the Credentials
// The channel can be used to retrieve messages
func (ttnClient *TTNClient) GetUplinkChannel() <-chan *types.UplinkMessage {
	// Subscribe to uplink messages
	uplink, err := ttnClient.devicePubSub.SubscribeUplink()
	if err != nil {
		log.Fatalf("%s: could not subscribe to uplink messages : %s", ttnClient.credentials.ClientName, err.Error())
	}
	return uplink
}

// Close ...
// Cleanup the mess we made
// Using defer would have been ideal but this gets the job done...
func (ttnClient *TTNClient) Close() {
	log.Println("Running TTNClient cleanup...")
	ttnClient.devicePubSub.Close()
	ttnClient.pubsub.Close()
	ttnClient.client.Close()
}
