// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"time"

	"github.com/platinasystems/goes"
	"github.com/platinasystems/goes-bmc/cmd/diag"
	"github.com/platinasystems/goes-bmc/cmd/fantrayd"
	"github.com/platinasystems/goes-bmc/cmd/fspd"
	"github.com/platinasystems/goes-bmc/cmd/ipcfg"
	"github.com/platinasystems/goes-bmc/cmd/ledgpiod"
	"github.com/platinasystems/goes-bmc/cmd/mmclog"
	"github.com/platinasystems/goes-bmc/cmd/mmclogd"
	"github.com/platinasystems/goes-bmc/cmd/qspi"
	"github.com/platinasystems/goes-bmc/cmd/toggle"
	"github.com/platinasystems/goes-bmc/cmd/ucd9090d"
	"github.com/platinasystems/goes-bmc/cmd/upgrade"
	"github.com/platinasystems/goes-bmc/cmd/w83795d"
	"github.com/platinasystems/goes/cmd"
	"github.com/platinasystems/goes/cmd/bang"
	"github.com/platinasystems/goes/cmd/buildid"
	"github.com/platinasystems/goes/cmd/buildinfo"
	"github.com/platinasystems/goes/cmd/cat"
	"github.com/platinasystems/goes/cmd/cd"
	"github.com/platinasystems/goes/cmd/chmod"
	"github.com/platinasystems/goes/cmd/cli"
	"github.com/platinasystems/goes/cmd/cmdline"
	"github.com/platinasystems/goes/cmd/cp"
	"github.com/platinasystems/goes/cmd/daemons"
	"github.com/platinasystems/goes/cmd/dhcpcd"
	"github.com/platinasystems/goes/cmd/dmesg"
	"github.com/platinasystems/goes/cmd/echo"
	eepromcmd "github.com/platinasystems/goes/cmd/eeprom"
	eeprom "github.com/platinasystems/goes/cmd/eeprom/platina_eeprom"
	"github.com/platinasystems/goes/cmd/elsecmd"
	"github.com/platinasystems/goes/cmd/env"
	"github.com/platinasystems/goes/cmd/exec"
	"github.com/platinasystems/goes/cmd/exit"
	"github.com/platinasystems/goes/cmd/export"
	"github.com/platinasystems/goes/cmd/falsecmd"
	"github.com/platinasystems/goes/cmd/femtocom"
	"github.com/platinasystems/goes/cmd/ficmd"
	"github.com/platinasystems/goes/cmd/flash_eraseall"
	"github.com/platinasystems/goes/cmd/for"
	"github.com/platinasystems/goes/cmd/function"
	"github.com/platinasystems/goes/cmd/gpio"
	"github.com/platinasystems/goes/cmd/grep"
	"github.com/platinasystems/goes/cmd/hdel"
	"github.com/platinasystems/goes/cmd/hdelta"
	"github.com/platinasystems/goes/cmd/hexists"
	"github.com/platinasystems/goes/cmd/hget"
	"github.com/platinasystems/goes/cmd/hgetall"
	"github.com/platinasystems/goes/cmd/hkeys"
	"github.com/platinasystems/goes/cmd/hset"
	"github.com/platinasystems/goes/cmd/i2c"
	"github.com/platinasystems/goes/cmd/i2cd"
	"github.com/platinasystems/goes/cmd/ifcmd"
	"github.com/platinasystems/goes/cmd/iminfo"
	"github.com/platinasystems/goes/cmd/imx6d"
	"github.com/platinasystems/goes/cmd/insmod"
	"github.com/platinasystems/goes/cmd/install"
	"github.com/platinasystems/goes/cmd/ip"
	"github.com/platinasystems/goes/cmd/kexec"
	"github.com/platinasystems/goes/cmd/keys"
	"github.com/platinasystems/goes/cmd/kill"
	"github.com/platinasystems/goes/cmd/ldp"
	"github.com/platinasystems/goes/cmd/ln"
	"github.com/platinasystems/goes/cmd/log"
	"github.com/platinasystems/goes/cmd/ls"
	"github.com/platinasystems/goes/cmd/lsmod"
	"github.com/platinasystems/goes/cmd/lsof"
	"github.com/platinasystems/goes/cmd/mkdir"
	"github.com/platinasystems/goes/cmd/mknod"
	"github.com/platinasystems/goes/cmd/mount"
	"github.com/platinasystems/goes/cmd/ping"
	"github.com/platinasystems/goes/cmd/ps"
	"github.com/platinasystems/goes/cmd/pwd"
	"github.com/platinasystems/goes/cmd/reboot"
	"github.com/platinasystems/goes/cmd/redisd"
	"github.com/platinasystems/goes/cmd/reload"
	"github.com/platinasystems/goes/cmd/restart"
	"github.com/platinasystems/goes/cmd/rm"
	"github.com/platinasystems/goes/cmd/rmmod"
	"github.com/platinasystems/goes/cmd/scp"
	"github.com/platinasystems/goes/cmd/slashinit"
	"github.com/platinasystems/goes/cmd/sleep"
	"github.com/platinasystems/goes/cmd/source"
	"github.com/platinasystems/goes/cmd/sshd"
	"github.com/platinasystems/goes/cmd/start"
	"github.com/platinasystems/goes/cmd/stop"
	"github.com/platinasystems/goes/cmd/stty"
	"github.com/platinasystems/goes/cmd/subscribe"
	"github.com/platinasystems/goes/cmd/sync"
	"github.com/platinasystems/goes/cmd/testcmd"
	"github.com/platinasystems/goes/cmd/thencmd"
	"github.com/platinasystems/goes/cmd/truecmd"
	"github.com/platinasystems/goes/cmd/ubi"
	"github.com/platinasystems/goes/cmd/umount"
	"github.com/platinasystems/goes/cmd/uptime"
	"github.com/platinasystems/goes/cmd/uptimed"
	"github.com/platinasystems/goes/cmd/version"
	"github.com/platinasystems/goes/cmd/watchdog"
	"github.com/platinasystems/goes/cmd/wget"
	"github.com/platinasystems/goes/cmd/while"

	"github.com/platinasystems/goes/external/redis/publisher"

	"github.com/platinasystems/goes/lang"
)

var Version = "(devel)"

var Goes = &goes.Goes{
	NAME: "goes-" + name,
	APROPOS: lang.Alt{
		lang.EnUS: "platina's mk1 baseboard management controller",
	},
	ByName: map[string]cmd.Cmd{
		"!":       bang.Command{},
		"cat":     cat.Command{},
		"cd":      &cd.Command{},
		"chmod":   chmod.Command{},
		"cli":     &cli.Command{},
		"cp":      cp.Command{},
		"daemons": daemons.Admin,
		"dhcpcd":  &dhcpcd.Command{},
		"diag":    diag.Command{},
		"dmesg":   dmesg.Command{},
		"echo":    echo.Command{},
		"eeprom":  eepromcmd.Command{},
		"else":    &elsecmd.Command{},
		"env":     &env.Command{},
		"exec":    exec.Command{},
		"exit":    exit.Command{},
		"export":  export.Command{},
		"false":   falsecmd.Command{},
		"fantrayd": &fantrayd.Command{
			Init: fantraydInit,
		},
		"femtocom":       femtocom.Command{},
		"fi":             &ficmd.Command{},
		"flash_eraseall": flash_eraseall.Command{},
		"for":            forcmd.Command{},
		"fspd": &fspd.Command{
			Init: fspdInit,
		},
		"function": &function.Command{},
		"goes-daemons": &daemons.Server{
			Init: [][]string{
				[]string{"redisd"},
				[]string{"fantrayd"},
				[]string{"fspd"},
				[]string{"i2cd"},
				[]string{"imx6d"},
				[]string{"ledgpiod"},
				[]string{"mmclogd"},
				[]string{"sshd"},
				[]string{"uptimed"},
				[]string{"ucd9090d"},
				[]string{"w83795d"},
				[]string{"watchdog"},
			},
		},
		"gpio":    gpio.Command{},
		"grep":    grep.Command{},
		"hdel":    hdel.Command{},
		"hdelta":  &hdelta.Command{},
		"hexists": hexists.Command{},
		"hget":    hget.Command{},
		"hgetall": hgetall.Command{},
		"hkeys":   hkeys.Command{},
		"hset":    hset.Command{},
		"i2c":     i2c.Command{},
		"i2cd":    i2cd.Command{},
		"if":      &ifcmd.Command{},
		"imx6d": &imx6d.Command{
			VpageByKey: map[string]uint8{
				"bmc.temperature.units.C": 1,
			},
		},
		"insmod": insmod.Command{},
		"install": &install.Command{
			DefaultArchive: "http://platina.io/goes/bmc/rootfs-arm.cpio.xz",
		},
		"ip":    ip.Goes,
		"ipcfg": ipcfg.Command{},
		"kexec": &kexec.Command{},
		"keys":  keys.Command{},
		"kill":  kill.Command{},
		"ldp":   ldp.Command{},
		"ledgpiod": &ledgpiod.Command{
			Init: ledgpiodInit,
		},
		"ln":      ln.Command{},
		"log":     log.Command{},
		"ls":      ls.Command{},
		"lsmod":   lsmod.Command{},
		"lsof":    lsof.Command{},
		"mkdir":   mkdir.Command{},
		"mknod":   mknod.Command{},
		"mmclog":  mmclog.Command{},
		"mmclogd": &mmclogd.Command{},
		"mount":   mount.Command{},
		"ping":    ping.Command{},
		"ps":      ps.Command{},
		"pwd":     pwd.Command{},
		"reboot":  &reboot.Command{},
		"redisd": &redisd.Command{
			Devs:    []string{"lo", "eth0"},
			Machine: "platina-mk1-bmc",
			Hook: func(pub *publisher.Publisher) {
				eeprom.Config(
					eeprom.BusIndex(0),
					eeprom.BusAddresses([]int{0x55}),
					eeprom.BusDelay(10*time.Millisecond),
					eeprom.MinMacs(2),
					eeprom.OUI([3]byte{0x02, 0x46, 0x8a}),
				)
				eeprom.RedisdHook(pub)
			},
		},
		"reload":  reload.Command{},
		"restart": &restart.Command{},
		"rm":      rm.Command{},
		"rmmod":   rmmod.Command{},
		"qspi":    qspi.Command{},
		"scp":     scp.Command{},
		"show": &goes.Goes{
			NAME:  "show",
			USAGE: "show OBJECT",
			APROPOS: lang.Alt{
				lang.EnUS: "print stuff",
			},
			ByName: map[string]cmd.Cmd{
				"buildid":   buildid.Command{},
				"buildinfo": buildinfo.Command{},
				"cmdline":   cmdline.Command{},
				"copyright": License,
				"iminfo":    iminfo.Command{},
				"license":   License,
				"log":       daemons.Log{},
				"machine":   Machine,
				"patents":   Patents,
				"version":   &version.Command{V: Version},
			},
		},
		"/init":  &slashinit.Command{FsHook: ubiSetup},
		"sleep":  sleep.Command{},
		"source": &source.Command{},
		"sshd":   &sshd.Command{FailSafe: false},
		"start": &start.Command{
			ConfGpioHook: startConfGpioHook,
			Gettys: []start.TtyCon{
				{Tty: "/dev/ttymxc0", Baud: 115200},
			},
		},
		"stop":      &stop.Command{},
		"stty":      stty.Command{},
		"subscribe": subscribe.Command{},
		"sync":      sync.Command{},
		"[":         testcmd.Command{},
		"then":      &thencmd.Command{},
		"toggle":    toggle.Command{},
		"true":      truecmd.Command{},
		"ubiattach": ubi.AttachCommand{},
		"ubidetach": ubi.DetachCommand{},
		"ucd9090d": &ucd9090d.Command{
			Init: ucd9090dInit,
		},
		"umount":  umount.Command{},
		"until":   whilecmd.Command{IsUntil: true},
		"upgrade": &upgrade.Command{},
		"uptime":  uptime.Command{},
		"uptimed": uptimed.Command{},
		"w83795d": &w83795d.Command{
			Init: w83795dInit,
		},
		"watchdog": &watchdog.Command{
			GpioPin: "BMC_WDI",
		},
		"wget":  wget.Command{},
		"while": whilecmd.Command{},
	},
}
