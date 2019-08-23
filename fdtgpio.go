// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/platinasystems/fdt"
	"github.com/platinasystems/gpio"
)

// Build map of gpio pins for this gpio controller
func GatherAliases(n *fdt.Node) {
	for p, pn := range n.Properties {
		if strings.Contains(p, "gpio") {
			val := strings.Split(string(pn), "\x00")
			v := strings.Split(val[0], "/")
			gpio.Aliases[p] = v[len(v)-1]
		}
	}
}

// Build map of gpio pins for this gpio controller
func GatherPins(n *fdt.Node, name string, value string) {
	var pn []string

	for na, al := range gpio.Aliases {
		if al == n.Name {
			for _, c := range n.Children {
				mode := ""
				for p, _ := range c.Properties {
					switch p {
					case "gpio-pin-desc":
						pn = strings.Split(c.Name, "@")
					case "output-high", "output-low", "input":
						mode = p
					}
				}
				err := gpio.NewPin(pn[0], mode, na, pn[1])
				if err != nil {
					fmt.Printf("Error setting %s to %s: %s\n",
						pn[0], mode, err)
				}
			}
		}
	}
}
