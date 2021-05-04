package main

import (
	"fmt"
)

func Example_ipCommand() {
	s := ipCommand("172.17.3.52::172.17.2.1:255.255.254.0::eth0")
	fmt.Println(string(s))
	// Output:
	// ip link eth0 change up
	// ip address add 172.17.3.52/23 dev eth0
	// ip route add 0.0.0.0/0 via 172.17.2.1
}
