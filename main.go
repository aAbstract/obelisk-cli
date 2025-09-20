package main

// device control
//	- read ltbus driver config from a json file
//	- connect to first avilable devices
//	- exec LTBus command

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aAbstract/obelisk_cli/lib"
)

func main() {
	fmt.Print(lib.ASCII_ART)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(lib.Get_cli_prompt())
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		root_cmd := strings.Fields(cmd)[0]

		switch root_cmd {
		case "list":
			lib.Obelisk_list(cmd)
		case "connect":
			lib.Obelisk_device_connect(cmd)
		case "disconnect":
			lib.Obelisk_device_disconnect(cmd)
		case "exec":
			lib.Obelisk_device_exec(cmd)
		default:
			fmt.Printf("Invalid Command: %s\n", cmd)
		}
	}
}
