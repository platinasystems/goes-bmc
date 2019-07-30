// Copyright © 2015-2017 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ipcfg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/platinasystems/goes/lang"
	"github.com/platinasystems/parms"
)

const (
	Machine = "platina-mk1-bmc"
)

type Command struct{}

func (Command) String() string { return "ipcfg" }

func (Command) Usage() string { return "ipcfg [-ip]" }

func (Command) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "set the bmc ip address, via bootargs ip=",
	}
}

func (Command) Man() lang.Alt {
	return lang.Alt{
		lang.EnUS: `
DESCRIPTION
        The ipcfg command sets the bmc legacy ip address.

		The legacy IP address is only used when installing BMC
		firmware releases which do not have persistent storage.

		The value is stored in the "ip" string in u-boot "bootargs"
		environment variable in legacy firmware versions.

		Example:
		ipcfg -ip 192.168.101.241::192.168.101.2:255.255.255.0::eth0:on`,
	}
}

func (Command) Main(args ...string) (err error) {
	parm, args := parms.New(args, "-ip")
	if len(parm.ByName["-ip"]) == 0 {
		if err = dispIP(); err != nil {
			return err
		}
		return
	} else {
		s := parm.ByName["-ip"]
		if err = updateIP(s); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func dispIP() error {
	per, err := ioutil.ReadFile("/boot/" + Machine + "-per.bin")
	if err != nil {
		return err
	}
	fmt.Println("ip=" + strings.TrimSuffix(string(per), "\x00"))

	return nil
}

func updateIP(ip string) (err error) {
	err = ioutil.WriteFile("/boot/"+Machine+"-per.bin",
		[]byte(ip+"\x00"), 0644)
	return
}
