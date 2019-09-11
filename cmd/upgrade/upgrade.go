// Copyright © 2015-2017 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

//TODO AUTHENTICATE,  CERT OR EQUIV
//TODO UPGRADE AUTOMATICALLY IF ENABLED, contact boot server

package upgrade

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/platinasystems/flags"
	"github.com/platinasystems/goes/lang"
	"github.com/platinasystems/parms"
	"github.com/platinasystems/ubi"
)

const (
	DfltMod     = 0755
	DfltSrv     = "downloads.platinasystems.com"
	DfltVer     = "LATEST"
	Machine     = "platina-mk1-bmc"
	ArchiveName = Machine + ".zip"
	VersionName = Machine + "-ver.bin"
	V2Name      = Machine + "-v2"
)

var legacy bool

type Command struct {
	Gpio func()
	gpio sync.Once
}

func (*Command) String() string { return "upgrade" }

func (*Command) Usage() string {
	return "upgrade [-v VER] [-s SERVER[/dir]] [-r] [-l] [-c] [-t] [-f]"
}

func (*Command) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "upgrade images",
	}
}

func (*Command) Man() lang.Alt {
	return lang.Alt{
		lang.EnUS: `
DESCRIPTION
	The upgrade command updates firmware images.

	The default upgrade version is "LATEST". 
	Or specify a version using "-v", in the form YYYYMMDD

	The -l flag display version of selected server and version.

	The -r flag reports QSPI version numbers and booted from.

	Images are downloaded from "downloads.platinasystems.com",
	Or from a server using "-s" followed by a URL or IPv4 address.

	Upgrade proceeds only if the selected version number is newer,
	unless overridden with the "-f" force flag.

OPTIONS
	-v [VER]          version [YYYYMMDD] or LATEST (default)
	-s [SERVER[/dir]] IP4 or URL, default downloads.platinasystems.com 
	-t                use TFTP instead of HTTP
	-l                display version of selected server and version
	-r                report QSPI installed version
	-c                check SHA-1's of flash
	-f                force upgrade (ignore version check)`,
	}
}

func (c *Command) Main(args ...string) error {
	flag, args := flags.New(args, "-t", "-l", "-f", "-r", "-c", "-legacy")
	parm, args := parms.New(args, "-v", "-s")
	if len(parm.ByName["-v"]) == 0 {
		parm.ByName["-v"] = DfltVer
	}
	if len(parm.ByName["-s"]) == 0 {
		parm.ByName["-s"] = DfltSrv
	}

	if flag.ByName["-l"] {
		if err := reportVerServer(parm.ByName["-s"], parm.ByName["-v"],
			flag.ByName["-t"]); err != nil {
			return err
		}
		return nil
	}
	if flag.ByName["-r"] {
		if err := reportVerQSPIdetail(); err != nil {
			return err
		}
		return nil
	}
	if flag.ByName["-c"] {
		if err := compareChecksums(); err != nil {
			return err
		}
		return nil
	}

	if err := doUpgrade(parm.ByName["-s"], parm.ByName["-v"],
		flag.ByName["-t"], flag.ByName["-f"], flag.ByName["-legacy"]); err != nil {
		return err
	}
	return nil
}

func reportVerServer(s string, v string, t bool) (err error) {
	n, err := getFile(s, v, t, ArchiveName)
	if err != nil {
		return fmt.Errorf("Error reading %s: %s\n", ArchiveName, err)
	}
	if n < 1000 {
		return fmt.Errorf("File %s %d bytes\n", ArchiveName, n)
	}
	if err := unzip(); err != nil {
		return fmt.Errorf("Server error: unzipping file: %\n", err)
	}
	defer rmFiles()

	l, err := ioutil.ReadFile(VersionName)
	if err != nil {
		fmt.Printf("Image version not found on server\n")
		return nil
	}
	sv := string(l[VERSION_OFFSET:VERSION_LEN])
	if string(l[VERSION_OFFSET:VERSION_DEV]) == "dev" {
		sv = "dev"
	}
	printVerServer(s, v, sv)
	return nil
}

func reportVerQSPIdetail() (err error) {
	err = printJSON()
	if err != nil {
		return err
	}
	return nil
}

func compareChecksums() (err error) {
	if err = cmpSums(); err != nil {
		return err
	}
	return nil
}

func doUpgrade(s string, v string, t bool, f bool, l bool) (err error) {
	n, err := getFile(s, v, t, ArchiveName)
	if err != nil {
		return fmt.Errorf("Error reading %s: %s\n", ArchiveName, err)
	}
	if n < 1000 {
		return fmt.Errorf("File %s %d bytes\n", ArchiveName, n)
	}
	if err = unzip(); err != nil {
		return fmt.Errorf("Server error: unzipping file: %v\n", err)
	}
	defer rmFiles()

	// If forced legacy, do that. Otherwise check if this is a UBI
	// volume. If not, always do legacy upgrades.
	isUbi, err := ubi.IsUbi(3)
	if err != nil {
		return err
	}
	if l || !isUbi {
		legacy = true
	} else {
		ubiMounted := true
		_, err = os.Stat("/sys/devices/virtual/ubi/ubi0")
		if err != nil {
			if os.IsNotExist(err) {
				ubiMounted = false
			} else {
				return fmt.Errorf("Unexpected error %s stating ubi device\n")
			}
		}

		_, err = os.Stat(V2Name)
		if l || os.IsNotExist(err) {
			if ubiMounted {
				return fmt.Errorf("Can't upgrade to legacy versions with UBI attached\n")
			}
			legacy = true
		} else {
			if err != nil {
				return fmt.Errorf("Unexpected error %s stating v2 file")
			}
			if !ubiMounted {
				return fmt.Errorf("Can't upgrade to current versions without UBI attached\n")
			}
		}
	}
	fmt.Print("\n")

	if !f {
		qv, err := getVerQSPI()
		if err != nil {
			return err
		}
		if len(qv) == 0 {
			fmt.Printf("Aborting, couldn't find version in QSPI\n")
			fmt.Printf("Use -f to force upgrade.\n")
			return nil
		}

		l, err := ioutil.ReadFile(VersionName)
		if err != nil {
			fmt.Printf("Aborting, couldn't find version number on server\n")
			fmt.Printf("Use -f to force upgrade.\n")
			return nil
		}
		sv := string(l[VERSION_OFFSET:VERSION_LEN])
		if string(l[VERSION_OFFSET:VERSION_DEV]) == "dev" {
			sv = "dev"
		}
		if sv != "dev" && qv != "dev" {
			newer, err := isVersionNewer(qv, sv)
			if err != nil {
				fmt.Printf("Aborting, server version error\n", sv)
				fmt.Printf("Use -f to force upgrade.\n")
				return nil
			}
			if !newer {
				fmt.Printf("Aborting, server version %s is not newer\n", sv)
				fmt.Printf("Use -f to force upgrade.\n")
				return nil
			}
		}
	}

	perFile, err := ioutil.ReadFile("/boot/" + Machine + "-per.bin")
	if err != nil {
		fmt.Printf("Error reading /boot/%s-per.bin: %s\n", Machine,
			err)
		fmt.Println("Using default of ip=dhcp")
		perFile = []byte("dhcp\x00")
	}
	err = ioutil.WriteFile(Machine+"-per.bin", perFile, 0644)
	if err != nil {
		return fmt.Errorf("Error creating %s-per.bin - aborting!")
	}
	if err = writeImageAll(); err != nil {
		return fmt.Errorf("*** UPGRADE ERROR! ***: %v\n", err)
	}
	UpdateEnv()

	return nil
}
