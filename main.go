package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"

	"github.com/j-forster/Wazihub-API"
)

var tThreshold float64 = 25

func main() {
	host.Init()

	log.Println("Wazihub Raspberry Pi Demo")

	//////////

	// The ID of the device this program is running on.
	// We use 'CurrentDeviceId' to create a unique Id for this device.
	log.Println("This device id is:", device.Id)

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

	go LoRa()
	go Bluetooth()

	// go ButtonListener()

	//////////

	gpio18 := gpioreg.ByName("18")
	gpio23 := gpioreg.ByName("23")

	// Let's hook into the actuation.
	led1, _ := wazihub.Actuation(device.Id, "led1")
	led2, _ := wazihub.Actuation(device.Id, "led2")
	tthreshold, _ := wazihub.Actuation(device.Id, "tthreshold")

	fmt.Println("Waiting for actuation...")

	for {
		select {
		case action := <-led1:
			log.Println("LED 1:", action)
			gpio18.Out(gpio.Level(action == "1"))

		case action := <-led2:
			log.Println("LED 2:", action)
			gpio23.Out(gpio.Level(action == "1"))

		case t := <-tthreshold:
			log.Println("Threshold:", t)
			tThreshold, _ = strconv.ParseFloat(t, 32)
		}
	}
}

//////////

// ButtonListener is a goroutine that waits for the switch button to change.
func ButtonListener() {

	gpio24 := gpioreg.ByName("24")
	// Connect a push-button to GPIO 24
	gpio24.In(gpio.PullUp, gpio.BothEdges)

	for {
		if gpio24.WaitForEdge(time.Second * 30) {
			wazihub.PostValue(device.Id, "push-button", gpio24.Read())
		}
	}
}

//////////

// looks like "T=26.19"
var bluetoothRegexp = regexp.MustCompile(`T=(\d*\.\d*)`)

func Bluetooth() {

	gpio18 := gpioreg.ByName("18")

	bluetooth, err := NewBluetoothDevice("E7:BC:F3:27:9A:CB")
	if err != nil {
		//bluetooth.Disconnect()
		log.Fatal(err)
	}
	err = bluetooth.Connect()
	if err != nil {
		//bluetooth.Disconnect()
		log.Fatal(err)
	}
	for {
		buf, err := bluetooth.Read("6e400003-b5a3-f393-e0a9-e50e24dcca9e")
		if err != nil {
			//bluetooth.Disconnect()
			log.Fatal(err)
		}
		log.Println("Bluetooth:", string(buf))
		values := bluetoothRegexp.FindStringSubmatch(string(buf))
		if values != nil {
			t, _ := strconv.ParseFloat(values[1], 32)
			wazihub.PostValue(device.Id, "t", t)

			if t > tThreshold {
				for i := 0; i < 5; i++ {
					gpio18.Out(gpio.Level(true))
					time.Sleep(time.Millisecond * 100)
					gpio18.Out(gpio.Level(false))
					time.Sleep(time.Millisecond * 100)
				}
				continue
			}
		}

		time.Sleep(time.Second)
	}
}

//////////

// looks like "SMC/427.70/BAT/16.0"
var loraRegexp = regexp.MustCompile(`SMC/(\d*\.\d*)/BAT/(\d*\.\d*)`)

func LoRa() {

	lora, err := NewLoRaDevice()
	if err != nil {
		log.Fatal(err)
	}
	for {
		time.Sleep(time.Second)
		buf, err := lora.Read()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("LoRa:", string(buf))
		values := loraRegexp.FindStringSubmatch(string(buf))
		if values != nil {
			smc, _ := strconv.ParseFloat(values[1], 32)
			bat, _ := strconv.ParseFloat(values[2], 32)
			wazihub.PostValue(device.Id, "smc", smc)
			wazihub.PostValue(device.Id, "bat", bat)
		}
	}
}

//////////
