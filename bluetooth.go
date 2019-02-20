package main

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type BluetoothDevice struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	reader *bufio.Reader
}

func NewBluetoothDevice(addr string) (*BluetoothDevice, error) {
	cmd := exec.Command("gatttool", "-t", "random", "-b", addr, "-I")
	device := &BluetoothDevice{cmd: cmd}
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

func (device *BluetoothDevice) Connect() error {
	device.stdin.Write([]byte("connect\n"))
	device.reader.ReadLine() // connect
	device.reader.ReadLine() // "Attempting to connect to ..."
	buff, _, err := device.reader.ReadLine()
	if err != nil {
		return err
	}
	line := string(buff)
	if !strings.HasSuffix(line, "Connection successful") {
		return errors.New(line)
	}
	return nil
}

func (device *BluetoothDevice) Read(uuid string) ([]byte, error) {
	device.stdin.Write([]byte("char-read-uuid " + uuid + "\n"))
	device.reader.ReadLine() // char-read-uuid
	buff, _, err := device.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	line := string(buff)
	if i := strings.Index(line, "value: "); i != -1 {
		chars := strings.Split(line[i+7:len(line)-2], " ")
		bytes := make([]byte, len(chars))
		for i, c := range chars {
			n, _ := strconv.ParseInt("0x"+c, 0, 64)
			bytes[i] = byte(n)
		}
		return bytes, nil
	}
	return nil, errors.New(line)
}

func (device *BluetoothDevice) Disconnect() error {
	device.stdin.Write([]byte("exit\n"))
	device.reader.ReadLine()
	return device.cmd.Wait()
}
