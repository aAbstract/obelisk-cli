package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DevicePort string `json:"device_port"`
	Driver     Driver `json:"driver"`
}

type Driver struct {
	DriverType   string       `json:"driver_type"`
	DriverConfig DriverConfig `json:"driver_config"`
}

type DriverConfig struct {
	LTBusMsgs []LTBusMsg `json:"ltbus_msgs"`
}

type LTBusMsg struct {
	RegisterName      string `json:"register_name"`
	RegisterAddr      string `json:"register_addr"`
	RegisterDataType  string `json:"register_data_type"`
	RegisterParamType int    `json:"register_param_type"`
}

var config *Config

func Get_conf() *Config {
	if config == nil {
		json_str, err := os.ReadFile("conf.json")
		if err != nil {
			fmt.Println("Error Reading File conf.json")
		}

		config = &Config{}
		err = json.Unmarshal(json_str, config)
		if err != nil {
			fmt.Println("Error Parsing conf.json")
		}
	}

	return config
}
