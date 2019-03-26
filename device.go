package main

import wazihub "github.com/j-forster/Wazihub-API"

var device = &wazihub.Device{
	Id:     wazihub.CurrentDeviceId(),
	Name:   "Raspberry Pi Test",
	Domain: "MyDomain",
	Sensors: []*wazihub.Sensor{
		&wazihub.Sensor{
			Id:            "Digital Soil Moisture",
			Name:          "SMC",
			SensingDevice: "",
			Unit:          "",
		},
		&wazihub.Sensor{
			Id:            "bat",
			Name:          "BAT",
			SensingDevice: "",
			Unit:          "",
		},
		&wazihub.Sensor{
			Id:            "t",
			Name:          "T",
			SensingDevice: "",
			Unit:          "°C",
		},
	},
	Actuators: []*wazihub.Actuator{
		&wazihub.Actuator{
			Id:              "tthreshold",
			Name:            "Alarm Temp. Threshold",
			ActuationDevice: "",
		},
		&wazihub.Actuator{
			Id:              "led1",
			Name:            "Status LED 1",
			ActuationDevice: "",
		},
		&wazihub.Actuator{
			Id:              "led2",
			Name:            "Status LED 2",
			ActuationDevice: "",
		},
	},
}
