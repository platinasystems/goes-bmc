// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"github.com/platinasystems/fdt"
	"github.com/platinasystems/gpio"
)

func gpioInit() {
	gpio.Aliases = make(gpio.GpioAliasMap)
	gpio.Pins = make(gpio.PinMap)

	t := fdt.DefaultTree()

	if t != nil {
		t.MatchNode("aliases", GatherAliases)
		t.EachProperty("gpio-controller", "", GatherPins)
	}
}
