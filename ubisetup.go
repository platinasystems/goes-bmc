// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/platinasystems/gpio"
	"github.com/platinasystems/mtd"
	"github.com/platinasystems/ubi"
)

type MTDinfo struct {
	typ       byte
	flags     uint32
	size      uint32
	erasesize uint32
	writesize uint32
	oobsize   uint32
	unused    uint64
}

type EraseInfo struct {
	start  uint32
	length uint32
}

const (
	MEMGETINFO = 0x80204d01
	MEMERASE   = 0x40084d02
)

// /etc => /perm/etc (recurse as /etc/ssl => /perm/etc/ssl
// /etc/ssl => /perm/etc/ssl
//
// Recursively calls itself as /etc
func copyRecurse(src, dst string, overwrite bool) (err error) {
	stat, err := os.Stat(src)
	if err != nil {
		return
	}
	if stat.Mode().IsDir() {
		files, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			err = os.Mkdir(dst, stat.Mode())
			if err != nil {
				return err
			}
			fmt.Printf("mkdir -m %o %s\n",
				stat.Mode()&os.ModePerm, dst)
		}
		for _, file := range files {
			base := path.Base(file.Name())
			err = copyRecurse(src+"/"+base, dst+"/"+base, overwrite)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if !stat.Mode().IsRegular() {
		return nil // Skip special files
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if _, err := os.Stat(dst); (overwrite && err == nil) ||
		os.IsNotExist(err) {

		destination, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
		fmt.Printf("copy %s=>%s\n", src, dst)
	}
	return nil
}

func ubiSetup() (err error) {
	pin, found := gpio.FindPin("QSPI_MUX_SEL")
	if found {
		r, _ := pin.Value()
		qs := 0
		dir := "low"
		if r {
			qs = 1
			dir = "high"
		}
		err = pin.SetDirection(dir)
		if err != nil {
			fmt.Printf("Error setting QSPI_MUX_SEL output: %s\n",
				err)
			return
		}
		fmt.Printf("Booted from QSPI%d\n", qs)
	} else {
		fmt.Print("Unable to find QSPI_MUX_SEL, skipping UBI setup")
		return
	}

	// Check if this is a UBI volume. If not, we need to convert it,
	// so stash the itb, per, and ver partitions.

	ubiDev, err := mtd.NameToUnit("ubi")
	if err != nil {
		return err
	}
	ubiDevName := "/dev/mtd" + strconv.Itoa(ubiDev)

	isUbi, err := ubi.IsUbi(int32(ubiDev))
	if err != nil {
		return err
	}
	if !isUbi {
		if _, err := os.Stat("/boot"); os.IsNotExist(err) {
			err = os.Mkdir("/boot", 0644)
			if err != nil {
				return err
			}
		}
		itbDev, err := mtd.NameToUnit("itb")
		if err != nil {
			return err
		}
		itbDevName := "/dev/mtd" + strconv.Itoa(itbDev)
		itb, err := ioutil.ReadFile(itbDevName)
		if err != nil {
			return fmt.Errorf("Error reading %s: %w", itbDevName, err)
		}
		itb = []byte(strings.TrimRight(string(itb), "\xff"))

		err = ioutil.WriteFile("/boot/platina-mk1-bmc-itb.bin",
			itb, 0644)
		if err != nil {
			return err
		}

		perDev, err := mtd.NameToUnit("per")
		if err != nil {
			return err
		}
		perDevName := "/dev/mtd" + strconv.Itoa(perDev)
		per, err := ioutil.ReadFile(perDevName)
		if err != nil {
			return fmt.Errorf("Error reading %s: %w", perDevName, err)
		}

		perNoNul := ""
		if nul := strings.Index(string(per), "\x00"); nul >= 0 {
			per = per[:nul+1]
			perNoNul = string(per[:nul])
		}
		err = ioutil.WriteFile("/boot/platina-mk1-bmc-per.bin",
			per, 0644)
		if err != nil {
			return err
		}
		ipCmd := ipCommand(perNoNul)
		if len(ipCmd) > 0 {
			err = ioutil.WriteFile("/etc/goes/start", ipCmd, 0644)
			if err != nil {
				return fmt.Errorf("Error writing /etc/goes/start: %s",
					err)
			}
		}
		verDev, err := mtd.NameToUnit("ver")
		if err != nil {
			return err
		}
		verDevName := "/dev/mtd" + strconv.Itoa(verDev)
		ver, err := ioutil.ReadFile(verDevName)
		if err != nil {
			return fmt.Errorf("Error reading %s: %w", verDevName, err)
		}
		ver = []byte(strings.TrimRight(string(ver), "\xff"))
		err = ioutil.WriteFile("/boot/platina-mk1-bmc-ver.bin",
			ver, 0644)
		if err != nil {
			return err
		}

		m, err := os.OpenFile(ubiDevName, os.O_RDWR, 0666)
		if err != nil {
			return fmt.Errorf("Unable to open %s: %w", ubiDevName, err)
		}
		mi := &MTDinfo{}
		_, _, e := syscall.Syscall(syscall.SYS_IOCTL, m.Fd(),
			uintptr(MEMGETINFO), uintptr(unsafe.Pointer(mi)))
		if e != 0 {
			return fmt.Errorf("Error getting info on %s: %w", ubiDevName, e)
		}
		ei := &EraseInfo{0, mi.erasesize}
		for ei.start = 0; ei.start < mi.size; ei.start += mi.erasesize {
			fmt.Println("Erasing Block...", ei.start, ei.length)
			_, _, e := syscall.Syscall(syscall.SYS_IOCTL, m.Fd(),
				uintptr(MEMERASE), uintptr(unsafe.Pointer(ei)))
			if e != 0 {
				fmt.Printf("Erase error block %d: %s", ei.start,
					e)
				return e
			}
		}
	}

	err = ubi.Attach(0, int32(ubiDev), 0, 0)
	if err != nil {
		return
	}

	_, err = ubi.FindVolume(0, "perm")
	if err == ubi.ErrVolumeNotFound {
		di, err := ubi.Info(0)
		if err != nil {
			return err
		}
		err = ubi.Mkvol(0, ubi.VolNumAuto, 1, false,
			int64(di.Avail_eraseblocks*di.Eraseblock_size), "perm")
		if err != nil {
			fmt.Printf("Error creating perm volume\n")
			return err
		}
	} else {
		if err != nil {
			fmt.Printf("Error finding perm volume: %s\n", err)
			return err
		}
	}

	if _, err := os.Stat("/perm"); os.IsNotExist(err) {
		err = os.Mkdir("/perm", os.FileMode(0666))
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("/volatile"); os.IsNotExist(err) {
		err = os.Mkdir("/volatile", os.FileMode(0666))
		if err != nil {
			return err
		}
	}

	perm, err := ubi.FindVolume(0, "perm")
	if err != nil {
		return err
	}

	fmt.Printf("Mounting /perm on %s\n", perm)
	err = syscall.Mount(perm, "/perm", "ubifs", uintptr(0), "")
	if err != nil {
		return err
	}

	for _, i := range []struct {
		dir       string
		overwrite bool
	}{
		{"/etc", false},
		{"/boot", true},
	} {
		err = copyRecurse(i.dir, "/perm"+i.dir, i.overwrite)
		if err != nil {
			return err
		}
		if _, err := os.Stat("/volatile" + i.dir); os.IsNotExist(err) {
			err = os.Mkdir("/volatile"+i.dir, os.FileMode(0666))
			if err != nil {
				return err
			}
		}
		err = syscall.Mount(i.dir, "/volatile"+i.dir, "",
			uintptr(syscall.MS_BIND), "")
		if err != nil {
			return err
		}
		err = syscall.Mount("/perm"+i.dir, i.dir, "",
			uintptr(syscall.MS_BIND), "")
		if err != nil {
			return err
		}
	}
	return
}

func ipCommand(per string) []byte {
	ipS := strings.Split(per, ":")
	outstr := ""

	if len(ipS) > 0 {
		clientIP := ipS[0]
		if clientIP == "dhcp" {
			outstr += "daemons start dhcpcd\n"
		} else {
			if len(ipS) > 3 {
				netmaskIP := net.ParseIP(ipS[3])
				if netmaskIP != nil {
					nm := netmaskIP.To4()
					if nm != nil {
						mask := net.IPv4Mask(nm[0], nm[1],
							nm[2], nm[3])
						ones, bits := mask.Size()
						if bits == net.IPv4len*8 {
							clientIP = clientIP + "/" +
								strconv.Itoa(ones)
						}
					}
				}
			}
			outstr += fmt.Sprintf(`ip link eth0 change up
ip address add %s dev eth0
`,
				clientIP)
		}
	}
	if len(ipS) > 2 {
		outstr += fmt.Sprintf("ip route add 0.0.0.0/0 via %s\n", ipS[2])
	}
	return []byte(outstr)
}
