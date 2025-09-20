package lib

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var device_service_channel chan string

func Get_device_service_channel() chan<- string {
	return device_service_channel
}

var cli_prompt string = "obelisk (disconnected)> "

func Get_cli_prompt() string {
	return cli_prompt
}

func dt_size(dt string) int {
	dt_size, _ := strconv.Atoi(dt[1:])
	return dt_size / 8
}

func exec_read_formated_register(serial_port *serial.Port, dcmd string) {
	dcmd_parts := strings.Fields(dcmd)
	addr, _ := strconv.ParseUint(dcmd_parts[1], 0, 16)
	data_type := dcmd_parts[2]
	size := dt_size(data_type)
	ltbus_req := LTBus_read_request(uint16(addr), uint16(size))
	serial_port.Write(ltbus_req)

	resp_size := 10 + size
	ltbus_resp := make([]byte, resp_size)
	n, err := serial_port.Read(ltbus_resp)
	if err != nil || n != resp_size {
		fmt.Println("Invalid LTBus Packet Size")
		return
	}
	if ltbus_resp[0] != 0x7B || ltbus_resp[resp_size-1] != 0x7D {
		fmt.Println("Invalid LTBus Packet Bounds")
		return
	}

	data_payload := ltbus_resp[LTBUS_DATA_START : LTBUS_DATA_START+size]
	switch data_type {
	case "U8":
		val := data_payload[0]
		fmt.Printf("LTBus Response: %X\n", val)
	case "U16":
		val := binary.LittleEndian.Uint16(data_payload)
		fmt.Printf("LTBus Response: %X\n", val)
	case "U32":
		val := binary.LittleEndian.Uint32(data_payload)
		fmt.Printf("LTBus Response: %X\n", val)
	case "U64":
		val := binary.LittleEndian.Uint64(data_payload)
		fmt.Printf("LTBus Response: %X\n", val)
	case "I8":
		val := int8(data_payload[0])
		fmt.Printf("LTBus Response: %d\n", val)
	case "I16":
		val := int16(binary.LittleEndian.Uint16(data_payload))
		fmt.Printf("LTBus Response: %d\n", val)
	case "I32":
		val := int32(binary.LittleEndian.Uint32(data_payload))
		fmt.Printf("LTBus Response: %d\n", val)
	case "I64":
		val := int64(binary.LittleEndian.Uint64(data_payload))
		fmt.Printf("LTBus Response: %d\n", val)
	case "F32":
		val := math.Float32frombits(binary.LittleEndian.Uint32(data_payload))
		fmt.Printf("LTBus Response: %f\n", val)
	case "F64":
		val := math.Float64frombits(binary.LittleEndian.Uint64(data_payload))
		fmt.Printf("LTBus Response: %f\n", val)
	}
}

func device_service(serial_port *serial.Port, _device_service_channel <-chan string) {
	defer serial_port.Close()

	for {
		dcmd := <-_device_service_channel
		if dcmd == "disconnect" {
			break
		}

		root_dcmd := strings.Fields(dcmd)[0]
		switch root_dcmd {
		case "RFR":
			exec_read_formated_register(serial_port, dcmd)
		}
	}

	fmt.Println("Device Service Stopped")
}

var is_device_connected bool = false

func Obelisk_device_connect(cmd string) {
	cmd_parts := strings.Fields(cmd)
	var device_port string
	if len(cmd_parts) == 1 {
		conf := Get_conf()
		device_port = conf.DevicePort
	} else if len(cmd_parts) == 2 {
		device_port = cmd_parts[1]
	} else {
		fmt.Printf("Invalid Command: %s\n", cmd)
	}

	fmt.Printf("Connecting to Device: %s...\n", device_port)
	conf := &serial.Config{Name: device_port, Baud: 115200, ReadTimeout: time.Second * 1}
	serial_port, err := serial.OpenPort(conf)
	if err != nil {
		fmt.Printf("Connecting to Device: %s...ERR\n", device_port)
		return
	}
	cli_prompt = fmt.Sprintf("obelisk (connected - %s)> ", device_port)
	is_device_connected = true
	fmt.Printf("Connecting to Device: %s...OK\n", device_port)
	device_service_channel = make(chan string)
	go device_service(serial_port, device_service_channel)
}

func Obelisk_device_disconnect(cmd string) {
	_device_service_channel := Get_device_service_channel()
	_device_service_channel <- "disconnect"
	is_device_connected = false
	cli_prompt = "obelisk (disconnected)> "
}

var exec_cmd_set = map[string]struct{}{
	"RR":  {},
	"RFR": {},
	"WR":  {},
}

var dt_set = map[string]struct{}{
	"U8":  {},
	"U16": {},
	"U32": {},
	"U64": {},

	"I8":  {},
	"I16": {},
	"I32": {},
	"I64": {},

	"F32": {},
	"F64": {},
}

func Obelisk_device_exec(cmd string) {
	if !is_device_connected {
		fmt.Println("No LTBus Device Connected")
		return
	}

	cmd_parts := strings.Fields(cmd)
	if len(cmd_parts) < 2 {
		fmt.Println("Missing device_exec Command")
		return
	}

	root_exec_cmd := cmd_parts[1]
	if _, ok := exec_cmd_set[root_exec_cmd]; !ok {
		fmt.Printf("Invalid device_exec Command: %s\n", root_exec_cmd)
		return
	}

	if root_exec_cmd == "RR" {
		fmt.Printf("device_exec Command `%s` not Implemented\n", root_exec_cmd)
		return
	}

	if root_exec_cmd == "RFR" {
		if len(cmd_parts) != 4 {
			fmt.Printf("Usage: RFR <address> <type> | Read and Parse Register from LT Bus")
			return
		}

		if _, err := strconv.ParseUint(cmd_parts[2], 0, 16); err != nil {
			fmt.Printf("Error Parsing LTBus Address: %s\n", cmd_parts[2])
			return
		}

		if _, ok := dt_set[cmd_parts[3]]; !ok {
			fmt.Printf("Invalid LTBus Data Type: %s\n", cmd_parts[3])
			return
		}

		_device_service_channel := Get_device_service_channel()
		_device_service_channel <- strings.Join(cmd_parts[1:], " ")
		return
	}
}
