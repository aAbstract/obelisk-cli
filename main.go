package main

// list connected devices

// device control
//	- read ltbus driver config from a json file
//	- connect to first avilable devices
//	- exec LTBus command

import (
	"log"
	"os"

	"github.com/aAbstract/obelisk_cli/lib"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid Command Line Arguments")
	}

	root_cmd := os.Args[1]
	if root_cmd == "list" {
		lib.Obelisk_list()
	} else {
		log.Fatalf("Invalid Command: %s\n", root_cmd)
	}
}
