package lib

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/tarm/serial"
)

func check_ltbus_device(port string) uint16 {
	conf := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second * 1}
	serial_port, err := serial.OpenPort(conf)
	if err != nil {
		return 0
	}

	read_req := LTBus_read_request(0xA000, 2)
	_, err = serial_port.Write(read_req)
	if err != nil {
		return 0
	}

	ltbus_resp := make([]byte, 12)
	n, err := serial_port.Read(ltbus_resp)
	if err != nil || n != 12 {
		return 0
	}
	if ltbus_resp[0] != 0x7B || ltbus_resp[11] != 0x7D {
		return 0
	}

	device_id := uint16(ltbus_resp[LTBUS_DATA_START+1])<<8 | uint16(ltbus_resp[LTBUS_DATA_START])
	return device_id
}

func list_devices_linux() {
	acm_devices, _ := filepath.Glob("/dev/ttyACM*")
	usb_devices, _ := filepath.Glob("/dev/ttyUSB*")
	devices := append(acm_devices, usb_devices...)

	fmt.Println("Scanning Devices...")
	for _, p := range devices {
		fmt.Printf("Scanning Port %s...\n", p)
		device_id := check_ltbus_device(p)
		if device_id != 0 {
			fmt.Printf("\tDetected LTBus Device Port: %s, Device_ID: 0x%X\n", p, device_id)
		}
	}
}

func list_devices_windows() {
	log.Fatal("Command `list` is not Implemented on Windows")
}

func list_devices_mac() {
	log.Fatal("Command `list` is not Implemented on MacOS")
}

func Obelisk_list() {
	switch runtime.GOOS {
	case "linux":
		list_devices_linux()
	case "windows":
		list_devices_windows()
	case "darwin":
		list_devices_mac()
	default:
		log.Fatalf("Command `list` -> Unsupported Platform: %s\n", runtime.GOOS)
	}
}
