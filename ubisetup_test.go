package main

import (
	"fmt"
	"testing"
)

func TestIpCommand(t *testing.T) {
	s := ipCommand([]byte("172.17.3.52::172.17.2.1:255.255.254.0::eth0"))
	fmt.Println(string(s))
}
