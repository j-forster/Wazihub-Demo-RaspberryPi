package main

import (
	"bufio"
	"io"
	"os/exec"
)

type LoRaDevice struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	reader *bufio.Reader
}

func NewLoRaDevice() (*LoRaDevice, error) {
	cmd := exec.Command("/home/pi/LoRaWAN/lora_gateway")
	device := &LoRaDevice{cmd: cmd}
	var err error
	device.stdin, err = cmd.StdinPipe()
	if err != nil {
		return device, err
	}
	device.stdout, err = cmd.StdoutPipe()
	if err != nil {
		return device, err
	}
	device.stderr, err = cmd.StderrPipe()
	if err != nil {
		return device, err
	}
	device.reader = bufio.NewReader(device.stdout)
	return device, cmd.Start()
}

func (device *LoRaDevice) Read() ([]byte, error) {
	buff, _, err := device.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	return buff[2:], nil
}
