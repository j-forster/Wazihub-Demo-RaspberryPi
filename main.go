package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"

	"github.com/j-forster/Wazihub-API"
)

var deviceId = wazihub.CurrentDeviceId()

func main() {
	host.Init()

	log.Println("Wazihub Raspberri Pi Demo")

	//////////

	// See 'device.json'.
	file, err := ioutil.ReadFile("device.json")
	if err != nil {
		log.Fatal(err)
	}
	// The ID of the device this program is running on.
	// We use 'CurrentDeviceId' to create a unique Id for this device.
	log.Println("This device id is:", deviceId)

	// Create a new Device ...
	device := &wazihub.Device{Id: deviceId}
	// ... with the data from 'device.json'
	err = json.Unmarshal(file, device)
	if err != nil {
		log.Fatal(err)
	}

	//////////

	// Login to Wazihub
	if err := wazihub.Login("cdupont", "password"); err != nil {
		log.Fatalln("Login failed!", err)
	}

	//////////

	// We register the device, even though it might already be registered ...
	err = wazihub.CreateDevice(device)
	if err != nil {
		log.Fatalln("Failed to register!", err)
	}
	log.Println("Device registered!")

	//////////

	go ButtonListener()

	//////////

	gpio18 := gpioreg.ByName("18")
	gpio23 := gpioreg.ByName("23")

	// Let's hook into the actuation.
	led1, _ := wazihub.Actuation(deviceId, "led1")
	led2, _ := wazihub.Actuation(deviceId, "led2")

	fmt.Println("Waiting for actuation...")

	for {
		select {
		case action := <-led1:
			log.Println("LED 1:", action)
			gpio18.Out(gpio.Level(action == "1"))

		case action := <-led2:
			log.Println("LED 2:", action)
			gpio23.Out(gpio.Level(action == "1"))
		}
	}
}

//////////

func ButtonListener() {

	gpio24 := gpioreg.ByName("24")
	// Connect at push-button to GPIO 24
	gpio24.In(gpio.PullUp, gpio.BothEdges)

	for {
		if gpio24.WaitForEdge(time.Second * 30) {
			wazihub.PostValue(deviceId, "push-button", gpio24.Read())
		}
	}
}
