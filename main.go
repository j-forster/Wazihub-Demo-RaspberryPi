package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/j-forster/Wazihub-API"
)

func main() {

	log.Println("Wazihub Raspberri Pi Demo")

	//////////

	// See 'device.json'.
	file, _ := ioutil.ReadFile("device.json")
	// The ID of the device this program is running on.
	// We use 'CurrentDeviceId' to create a unique Id for this device.
	deviceId := wazihub.CurrentDeviceId()
	log.Println("This device id is:", deviceId)

	// Create a new Device ...
	device := &wazihub.Device{Id: deviceId}
	// ... with the data from 'device.json'
	json.Unmarshal(file, device)

	//////////

	// Login to Wazihub
	if err := wazihub.Login("cdupont", "password"); err != nil {
		log.Fatalln("Login failed!", err)
	}

	//////////

	// We register the device, even though it might already be registered ...
	err := wazihub.CreateDevice(device)
	if err != nil {
		log.Fatalln("Failed to register!", err)
	}
	log.Println("Device registered!")

	//////////

	ioutil.WriteFile("/sys/class/gpio/export", []byte("18"), 0666)
	ioutil.WriteFile("/sys/class/gpio/gpio18/direction", []byte("out"), 0666)

	ioutil.WriteFile("/sys/class/gpio/export", []byte("23"), 0666)
	ioutil.WriteFile("/sys/class/gpio/gpio23/direction", []byte("out"), 0666)

	//////////

	// Let's hook into the actuation.
	led1, _ := wazihub.Actuation(deviceId, "led1")
	led2, _ := wazihub.Actuation(deviceId, "led2")

	fmt.Println("Waiting for actuation...")

	for {
		select {
		case action := <-led1:
			log.Println("LED 1:", action)
			ioutil.WriteFile("/sys/class/gpio/gpio18/value", []byte(action), 0666)

		case action := <-led2:
			log.Println("LED 2:", action)
			ioutil.WriteFile("/sys/class/gpio/gpio23/value", []byte(action), 0666)
		}
	}
}
