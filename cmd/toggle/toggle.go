// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package toggle

import (
	"time"

	"github.com/platinasystems/goes/lang"
	"github.com/platinasystems/gpio"
	"github.com/platinasystems/i2c"
)

const i2cGpioAddr = 0x74

type Command struct{}

func (Command) String() string { return "toggle" }

func (Command) Usage() string { return "toggle SECONDS" }

func (Command) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "toggle console port between x86 and BMC",
	}
}

func (Command) Man() lang.Alt {
	return lang.Alt{
		lang.EnUS: `
DESCRIPTION
	The toggle command toggles the console port between x86 and BMC.`,
	}
}

func (Command) Main(args ...string) error {
	pin, found := gpio.FindPin("CPU_TO_MAIN_I2C_EN")
	if found {
		pin.SetValue(true)
	}
	time.Sleep(10 * time.Millisecond)
	uartToggle()
	if found {
		pin.SetValue(false)
	}
	pin, found = gpio.FindPin("FP_BTN_UARTSEL_EN_L")
	if found {
		pin.SetValue(true)
	}
	time.Sleep(10 * time.Millisecond)

	return nil
}

func uartToggle() {
	var dir0, out0 uint8
	i2c.Do(0, i2cGpioAddr,
		func(bus *i2c.Bus) (err error) {
			var d i2c.SMBusData
			reg := uint8(6)
			err = bus.Read(reg, i2c.WordData, &d)
			dir0 = d[0]
			return
		})
	i2c.Do(0, i2cGpioAddr,
		func(bus *i2c.Bus) (err error) {
			var d i2c.SMBusData
			reg := uint8(6)
			err = bus.Read(reg, i2c.WordData, &d)
			out0 = d[0]
			return
		})
	i2c.Do(0, i2cGpioAddr,
		func(bus *i2c.Bus) (err error) {
			var d i2c.SMBusData
			d[0] = out0 | 0x20
			reg := uint8(2)
			err = bus.Write(reg, i2c.ByteData, &d)
			return
		})
	i2c.Do(0, i2cGpioAddr,
		func(bus *i2c.Bus) (err error) {
			var d i2c.SMBusData
			d[0] = dir0 ^ 0x20
			reg := uint8(6)
			err = bus.Write(reg, i2c.ByteData, &d)
			return
		})

}
