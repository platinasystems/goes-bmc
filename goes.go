// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"github.com/platinasystems/go/internal/goes"
	"github.com/platinasystems/go/internal/goes/cmd/apropos"
	"github.com/platinasystems/go/internal/goes/cmd/bang"
	"github.com/platinasystems/go/internal/goes/cmd/boot"
	"github.com/platinasystems/go/internal/goes/cmd/cat"
	"github.com/platinasystems/go/internal/goes/cmd/cd"
	"github.com/platinasystems/go/internal/goes/cmd/chmod"
	"github.com/platinasystems/go/internal/goes/cmd/cli"
	"github.com/platinasystems/go/internal/goes/cmd/cmdline"
	"github.com/platinasystems/go/internal/goes/cmd/complete"
	"github.com/platinasystems/go/internal/goes/cmd/cp"
	"github.com/platinasystems/go/internal/goes/cmd/daemons"
	"github.com/platinasystems/go/internal/goes/cmd/dmesg"
	"github.com/platinasystems/go/internal/goes/cmd/echo"
	"github.com/platinasystems/go/internal/goes/cmd/eeprom"
	"github.com/platinasystems/go/internal/goes/cmd/env"
	"github.com/platinasystems/go/internal/goes/cmd/exec"
	"github.com/platinasystems/go/internal/goes/cmd/exit"
	"github.com/platinasystems/go/internal/goes/cmd/export"
	"github.com/platinasystems/go/internal/goes/cmd/fantray"
	"github.com/platinasystems/go/internal/goes/cmd/femtocom"
	"github.com/platinasystems/go/internal/goes/cmd/fsp"
	"github.com/platinasystems/go/internal/goes/cmd/gpio"
	"github.com/platinasystems/go/internal/goes/cmd/hdel"
	"github.com/platinasystems/go/internal/goes/cmd/hdelta"
	"github.com/platinasystems/go/internal/goes/cmd/help"
	"github.com/platinasystems/go/internal/goes/cmd/hexists"
	"github.com/platinasystems/go/internal/goes/cmd/hget"
	"github.com/platinasystems/go/internal/goes/cmd/hgetall"
	"github.com/platinasystems/go/internal/goes/cmd/hkeys"
	"github.com/platinasystems/go/internal/goes/cmd/hset"
	"github.com/platinasystems/go/internal/goes/cmd/i2c"
	"github.com/platinasystems/go/internal/goes/cmd/i2cd"
	"github.com/platinasystems/go/internal/goes/cmd/iminfo"
	"github.com/platinasystems/go/internal/goes/cmd/imx6"
	"github.com/platinasystems/go/internal/goes/cmd/insmod"
	"github.com/platinasystems/go/internal/goes/cmd/install"
	"github.com/platinasystems/go/internal/goes/cmd/kexec"
	"github.com/platinasystems/go/internal/goes/cmd/keys"
	"github.com/platinasystems/go/internal/goes/cmd/kill"
	"github.com/platinasystems/go/internal/goes/cmd/license"
	"github.com/platinasystems/go/internal/goes/cmd/ln"
	"github.com/platinasystems/go/internal/goes/cmd/log"
	"github.com/platinasystems/go/internal/goes/cmd/ls"
	"github.com/platinasystems/go/internal/goes/cmd/lsmod"
	"github.com/platinasystems/go/internal/goes/cmd/man"
	"github.com/platinasystems/go/internal/goes/cmd/mkdir"
	"github.com/platinasystems/go/internal/goes/cmd/mknod"
	"github.com/platinasystems/go/internal/goes/cmd/mount"
	"github.com/platinasystems/go/internal/goes/cmd/nlcounters"
	"github.com/platinasystems/go/internal/goes/cmd/nld"
	"github.com/platinasystems/go/internal/goes/cmd/nldump"
	"github.com/platinasystems/go/internal/goes/cmd/nsid"
	"github.com/platinasystems/go/internal/goes/cmd/patents"
	"github.com/platinasystems/go/internal/goes/cmd/ping"
	"github.com/platinasystems/go/internal/goes/cmd/platina/mk1/bmc/diag"
	"github.com/platinasystems/go/internal/goes/cmd/platina/mk1/bmc/ledgpio"
	"github.com/platinasystems/go/internal/goes/cmd/platina/mk1/toggle"
	"github.com/platinasystems/go/internal/goes/cmd/ps"
	"github.com/platinasystems/go/internal/goes/cmd/pwd"
	"github.com/platinasystems/go/internal/goes/cmd/reboot"
	"github.com/platinasystems/go/internal/goes/cmd/redisd"
	"github.com/platinasystems/go/internal/goes/cmd/reload"
	"github.com/platinasystems/go/internal/goes/cmd/resize"
	"github.com/platinasystems/go/internal/goes/cmd/restart"
	"github.com/platinasystems/go/internal/goes/cmd/rm"
	"github.com/platinasystems/go/internal/goes/cmd/rmmod"
	"github.com/platinasystems/go/internal/goes/cmd/show_commands"
	"github.com/platinasystems/go/internal/goes/cmd/slashinit"
	"github.com/platinasystems/go/internal/goes/cmd/sleep"
	"github.com/platinasystems/go/internal/goes/cmd/source"
	"github.com/platinasystems/go/internal/goes/cmd/start"
	"github.com/platinasystems/go/internal/goes/cmd/stop"
	"github.com/platinasystems/go/internal/goes/cmd/stty"
	"github.com/platinasystems/go/internal/goes/cmd/subscribe"
	"github.com/platinasystems/go/internal/goes/cmd/sync"
	"github.com/platinasystems/go/internal/goes/cmd/telnetd"
	"github.com/platinasystems/go/internal/goes/cmd/ucd9090"
	"github.com/platinasystems/go/internal/goes/cmd/umount"
	"github.com/platinasystems/go/internal/goes/cmd/uninstall"
	"github.com/platinasystems/go/internal/goes/cmd/uptimed"
	"github.com/platinasystems/go/internal/goes/cmd/usage"
	"github.com/platinasystems/go/internal/goes/cmd/version"
	"github.com/platinasystems/go/internal/goes/cmd/w83795"
	"github.com/platinasystems/go/internal/goes/cmd/watchdog"
	"github.com/platinasystems/go/internal/goes/cmd/wget"
)

var Goes goes.ByName

func init() {
	Goes = make(goes.ByName)
	Goes.Plot(apropos.New(),
		bang.New(),
		boot.New(),
		cat.New(),
		cd.New(),
		chmod.New(),
		cli.New(),
		cmdline.New(),
		complete.New(),
		cp.New(),
		daemons.New(),
		diag.New(),
		dmesg.New(),
		echo.New(),
		eeprom.New(),
		env.New(),
		exec.New(),
		exit.New(),
		export.New(),
		fantray.New(),
		femtocom.New(),
		fsp.New(),
		gpio.New(),
		hdel.New(),
		hdelta.New(),
		help.New(),
		hexists.New(),
		hget.New(),
		hgetall.New(),
		hkeys.New(),
		hset.New(),
		i2c.New(),
		i2cd.New(),
		iminfo.New(),
		imx6.New(),
		insmod.New(),
		install.New(),
		kexec.New(),
		keys.New(),
		kill.New(),
		ledgpio.New(),
		license.New(),
		ln.New(),
		log.New(),
		ls.New(),
		lsmod.New(),
		man.New(),
		mkdir.New(),
		mknod.New(),
		mount.New(),
		nlcounters.New(),
		nld.New(),
		nldump.New(),
		nsid.New(),
		patents.New(),
		ping.New(),
		ps.New(),
		pwd.New(),
		reboot.New(),
		redisd.New(),
		reload.New(),
		resize.New(),
		restart.New(),
		rm.New(),
		rmmod.New(),
		show_commands.New(),
		slashinit.New(),
		sleep.New(),
		source.New(),
		start.New(),
		stop.New(),
		stty.New(),
		subscribe.New(),
		sync.New(),
		telnetd.New(),
		toggle.New(),
		ucd9090.New(),
		umount.New(),
		uninstall.New(),
		uptimed.New(),
		usage.New(),
		version.New(),
		w83795.New(),
		watchdog.New(),
		wget.New(),
	)
}
