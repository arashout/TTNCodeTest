# TTNCodeTest

# Build
In your terminal run:
```
go build
.\TTNCodeTest.exe
```

# Setup 
Note that you need a .devenv.json file with your credentials which looks like this, with the fields filled in:
```
{
	"ttnAppId":"",
	"ttnAppAccessKey":"",
    "ttnDeviceId":"",
	"openSensorAPIKey":"",
	"openSensorClientId":"",
	"openSensorClientPassword":"",
	"openSensorTopicUrl":""
}
```
The `openSensorTopicUrl` should look like: `https://realtime.opensensors.io/v1/topics//users/username/homeoffice`