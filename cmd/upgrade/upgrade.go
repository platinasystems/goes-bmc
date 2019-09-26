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
	"path/filepath"

	"github.com/platinasystems/flags"
	"github.com/platinasystems/goes"
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
	TmpDir      = "/var/run/goes/upgrade"
)

var legacy bool

type Command struct {
	g *goes.Goes
}

func (c *Command) Goes(g *goes.Goes) { c.g = g }

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
	-f                force upgrade (ignore version check)
	-legacy           install legacy version`,
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

	isUbi, err := ubi.IsUbi(3)
	if err != nil {
		return fmt.Errorf("Error determining if UBI format: %s", err)
	}
	isUbiMounted, err := ubi.IsUbiMounted(0, 0)
	if err != nil {
		return fmt.Errorf("error determining if UBI mounted: %s", err)
	}
	mpc := 0 // Mount Point Count
	for _, mp := range []string{
		"/perm",
		"/boot",
		"/etc",
	} {
		dev, err := findDevForMountpoint(mp)
		if err != nil {
			return fmt.Errorf("Unable to determine where %s is mounted: %s",
				mp, err)
		}
		if dev == "" {
			continue
		}
		if dev != "/dev/ubi0_0" {
			return fmt.Errorf("Unexpected mount on %s: %s",
				mp, dev)
		}
		mpc++
	}

	if !isUbi && (isUbiMounted || mpc > 0) {
		return fmt.Errorf("Not UBI but mounted - internal error!")
	}
	if isUbi && (!isUbiMounted || mpc != 3) {
		return fmt.Errorf("Can't update UBI versions when unmounted, do qspi -mount")
	}
	err = os.MkdirAll(TmpDir, DfltMod)
	if err != nil {
		return fmt.Errorf("Unable to create work directory: %s", err)
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

	if err := c.doUpgrade(isUbi, parm.ByName["-s"], parm.ByName["-v"],
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

	l, err := ioutil.ReadFile(filepath.Join(TmpDir, VersionName))
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

func (c *Command) doUpgrade(isUbi bool, s string, v string, t bool, f bool, l bool) (err error) {
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

	if l || !isUbi {
		legacy = true
	} else {
		if _, err = os.Stat(filepath.Join(TmpDir, V2Name)); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Must use -legacy option to downgrade to legacy versions")
			}
			return fmt.Errorf("Unexpected error %s stating v2 file",
				err)
		}
	}

	if !f {
		qv, err := GetVerArchive()
		if err != nil {
			return err
		}
		if len(qv) == 0 {
			fmt.Printf("Aborting, couldn't find version in QSPI\n")
			fmt.Printf("Use -f to force upgrade.\n")
			return nil
		}

		l, err := ioutil.ReadFile(filepath.Join(TmpDir, VersionName))
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
				fmt.Printf("Aborting, server version %s is older than %s\n",
					sv, qv)
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
	err = ioutil.WriteFile(filepath.Join(TmpDir, Machine+"-per.bin"), perFile,
		0644)
	if err != nil {
		return fmt.Errorf("Error creating %s-per.bin - aborting!")
	}

	// If we are explicitly forcing legacy, and downgrading from UBI
	// to pre-UBI, we must unmount/detach for the user here. This is
	// because the UBI has to be attached for version checks, etc
	// to work correcly and so the user has to keep things mounted,
	// but we do the unmount just before we program the image.
	if isUbi && l {
		err = c.g.Main("qspi", "-unmount")
		if err != nil {
			return fmt.Errorf("Error unmounting UBI for legacy downgrade: %s",
				err)
		}
	}

	if err = writeImageAll(); err != nil {
		return fmt.Errorf("*** UPGRADE ERROR! ***: %v\n", err)
	}
	UpdateEnv()

	return nil
}
