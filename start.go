// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/platinasystems/goes/external/redis"
	"github.com/platinasystems/gpio"
	"github.com/platinasystems/log"
)

const (
	MaxEpollEvents = 32
	KB             = 1024
)

func startConfGpioHook() error {
	var deviceVer byte

	pin, found := gpio.FindPin("QSPI_MUX_SEL")
	if found {
		r, _ := pin.Value()
		if r {
			log.Print("Booted from QSPI1")
		} else {
			log.Print("Booted from QSPI0")
		}

	}

	for name, pin := range gpio.AllPins() {
		if name == "QSPI_MUX_SEL" {
			continue
		}
		err := pin.SetDefault()
		if err != nil {
			fmt.Printf("%s: %v\n", name, err)
		}
	}
	pin, found = gpio.FindPin("FRU_I2C_MUX_RST_L")
	if found {
		pin.SetValue(false)
		time.Sleep(1 * time.Microsecond)
		pin.SetValue(true)
	}

	pin, found = gpio.FindPin("MAIN_I2C_MUX_RST_L")
	if found {
		pin.SetValue(false)
		time.Sleep(1 * time.Microsecond)
		pin.SetValue(true)
	}

	redis.Hwait(redis.DefaultHash, "redis.ready", "true",
		10*time.Second)

	ss, _ := redis.Hget(redis.DefaultHash, "eeprom.DeviceVersion")
	_, _ = fmt.Sscan(ss, &deviceVer)
	if deviceVer == 0x0 || deviceVer == 0xff {
		pin, found = gpio.FindPin("FP_BTN_UARTSEL_EN_L")
		if found {
			pin.SetValue(false)
		}
	}
	return nil
}
