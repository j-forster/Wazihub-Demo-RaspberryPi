package main

import wazihub "github.com/j-forster/Wazihub-API"

var device = &wazihub.Device{
	Id:     wazihub.CurrentDeviceId(),
	Name:   "Raspberry Pi Test",
	Domain: "MyDomain",
	Sensors: []*wazihub.Sensor{
		&wazihub.Sensor{
			Id:            "push-button",
			Name:          "Push Button",
			SensingDevice: "PushButton",
			Unit:          "",
		},
	},
	Actuators: []*wazihub.Actuator{
		&wazihub.Actuator{
			Id:              "led1",
			Name:            "LED 1.",
			ActuationDevice: "lamp",
			ControlType:     "bool",
			Value:           false,
		},
		&wazihub.Actuator{
			Id:              "led2",
			Name:            "LED 2.",
			ActuationDevice: "lamp",
			ControlType:     "bool",
			Value:           false,
		},
	},
}
