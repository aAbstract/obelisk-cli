package lib_test

import (
	"slices"
	"testing"

	"github.com/aAbstract/obelisk_cli/lib"
)

func Test_ltbus_read_request(t *testing.T) {
	// RFR 0xD004 F32 -> 0x7B 0x01 0xAA 0x04 0xD0 0x04 0x00 0x7A 0xD3 0x7D
	request_packet := lib.LTBus_read_request(0xD004, 4)
	expected_packet := []byte{0x7B, 0x01, 0xAA, 0x04, 0xD0, 0x04, 0x00, 0x7A, 0xD3, 0x7D}
	if !slices.Equal(request_packet, expected_packet) {
		t.Error("Invalid Request Packet")
	}
}
