// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package upgrade

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/platinasystems/i2c"
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
	MEMGETINFO     = 0x80204d01 //from linux: mtd-abi.h
	MEMERASE       = 0x40084d02
	MEMLOCK        = 0x40084d05
	MEMUNLOCK      = 0x40084d06
	MEMERASE64     = 0x40104d14
	MTDdevice      = "/dev/mtd0"
	VERSION_OFFSET = 0x000
	VERSION_LEN    = 0x008
	VERSION_DEV    = 0x003
	JSON_OFFSET    = 0x100
	ENVSIZE        = 8192
	ENVCRC         = 4
)

type FlashFmt struct {
	off uint32
	siz uint32
}

var qFmt = map[string]FlashFmt{
	"ubo": FlashFmt{off: 0x000000, siz: 0x080000},
	"dtb": FlashFmt{off: 0x080000, siz: 0x040000},
	"env": FlashFmt{off: 0x0c0000, siz: 0x040000},
	"itb": FlashFmt{off: 0x100000, siz: 0x800000},
	"ker": FlashFmt{off: 0x100000, siz: 0x200000},
	"ini": FlashFmt{off: 0x300000, siz: 0x300000},
	"per": FlashFmt{off: 0xf80000, siz: 0x040000},
	"ver": FlashFmt{off: 0xfc0000, siz: 0x040000},
}

var img = []string{"ubo", "dtb", "env"}
var legacyImg = []string{"ker", "ini", "itb", "per", "ver"}
var newImg = []string{"itb", "per", "ver"}

var mi = &MTDinfo{0, 0, 0, 0, 0, 0, 0}
var ei = &EraseInfo{0, 0}
var fd int = 0

var sd i2c.SMBusData

func readBlk(blknam string) (b []byte, err error) {
	_, b, err = readFlash(qFmt[blknam].off, qFmt[blknam].siz)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readFlash(of uint32, sz uint32) (n int, b []byte, err error) {
	fd, err = syscall.Open(MTDdevice, syscall.O_RDWR, 0)
	if err != nil {
		err = fmt.Errorf("Open error %s: %s", MTDdevice, err)
		return 0, nil, err
	}
	defer syscall.Close(fd)

	if err = infoQSPI(); err != nil {
		return 0, nil, err
	}
	if n, b, err = readQSPI(of, sz); err != nil {
		return 0, nil, err
	}
	return n, b, nil
}

func writeBlk(blknam string, b []byte) (err error) {
	size := qFmt[blknam].siz
	ba := make([]byte, size)
	for i, _ := range ba {
		ba[i] = 0xff
	}
	for i, _ := range b {
		ba[i] = b[i]
	}
	err = writeArrayVerify(ba, qFmt[blknam].off, qFmt[blknam].siz, true)
	if err != nil {
		return err
	}
	return nil
}

func writeImageAll() (err error) {
	fd, err = syscall.Open(MTDdevice, syscall.O_RDWR, 0)
	if err != nil {
		err = fmt.Errorf("Open error %s: %s", MTDdevice, err)
		return err
	}
	defer syscall.Close(fd)

	if err = infoQSPI(); err != nil {
		return err
	}
	for _, j := range img {
		if err := writeImageVerify(Machine+"-"+j+".bin",
			qFmt[j].off, qFmt[j].siz, true); err != nil {
			return err
		}
	}
	if !legacy {
		for _, j := range newImg {
			src := Machine + "-" + j + ".bin"
			s, err := ioutil.ReadFile(filepath.Join(TmpDir, src))
			if err != nil {
				return fmt.Errorf("Error reading %s: %s\n", src, err)
			}
			dst := "/boot/" + src
			err = os.Remove(dst)
			if err != nil {
				fmt.Printf("Error removing %s: %s\n",
					src, err)
			}
			err = ioutil.WriteFile(dst, s, 0644)
			if err != nil {
				return fmt.Errorf("Error writing %s: %s\n", dst, err)
			}
			fmt.Printf("Installed %s\n", dst)
		}
	} else {
		for _, j := range legacyImg {
			if err := writeImageVerify(Machine+"-"+j+".bin",
				qFmt[j].off, qFmt[j].siz, true); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeImageVerify(imBase string, of uint32, sz uint32, vf bool) error {
	im := filepath.Join(TmpDir, imBase)
	if fi, err := os.Stat(im); !os.IsNotExist(err) {
		if fi.Size() == 0 {
			fmt.Println("skipping file...", im)
			return nil
		}
		if err = eraseQSPI(of, sz); err != nil {
			return err
		}
		b, err := ioutil.ReadFile(im)
		if err != nil {
			return err
		}
		fmt.Println("Programming...", im)
		_, err = writeQSPI(b, of)
		if err != nil {
			return err
		}
		if vf {
			nn, bb, err := readQSPI(of, sz)
			if err != nil {
				err = fmt.Errorf("Read error: %s %v",
					im, err)
				return err
			}
			if nn != int(sz) {
				err = fmt.Errorf("Size error %v!=%v: %s %v",
					nn, sz, im, err)
				return err
			}
			//TODO REPLACE WITH DEEP EQUAL
			for i := range b {
				if b[i] != bb[i] {
					err = fmt.Errorf("Verify error: %s %v",
						im, err)
					return err
				}
			}
			fmt.Println("Verify passed:", im)
		}
	}
	return nil
}

func writeArrayVerify(b []byte, of uint32, sz uint32, vf bool) (err error) {
	fd, err = syscall.Open(MTDdevice, syscall.O_RDWR, 0)
	if err != nil {
		err = fmt.Errorf("Open error %s: %s", MTDdevice, err)
		return err
	}
	defer syscall.Close(fd)

	if err = infoQSPI(); err != nil {
		return err
	}
	if len(b) != int(sz) {
		err = fmt.Errorf("Array size doesn't match")
		return err
	}
	if err = eraseQSPI(of, sz); err != nil {
		return err
	}
	fmt.Println("Programming...")
	_, err = writeQSPI(b, of)
	if err != nil {
		return err
	}
	if vf {
		nn, bb, err := readQSPI(of, sz)
		if err != nil {
			err = fmt.Errorf("Read error: %v", err)
			return err
		}
		if nn != int(sz) {
			err = fmt.Errorf("Size error %v!=%v: %v",
				nn, sz, err)
			return err
		}
		//TODO REPLACE WITH DEEP EQUAL
		for i := range b {
			if b[i] != bb[i] {
				err = fmt.Errorf("Verify error: %v", err)
				return err
			}
		}
		fmt.Println("Verify passed")
	}
	return nil
}

func infoQSPI() (err error) {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(MEMGETINFO), uintptr(unsafe.Pointer(mi)))
	if e != 0 {
		err = fmt.Errorf("Open error : %s", e)
		return err
	}
	return nil
}

func readQSPI(of uint32, sz uint32) (int, []byte, error) {
	b := make([]byte, sz)
	_, err := syscall.Seek(fd, int64(of), 0)
	if err != nil {
		err = fmt.Errorf("Seek error: %x: %s", of, err)
		return 0, b, err
	}
	n, err := syscall.Read(fd, b)
	if err != nil {
		err = fmt.Errorf("Read error %x: %s", of, err)
		return 0, b, err
	}
	return n, b, nil
}

func writeQSPI(b []byte, of uint32) (int, error) {
	_, err := syscall.Seek(fd, int64(of), 0)
	if err != nil {
		err = fmt.Errorf("Seek error: %d: %s", of, err)
		return 0, err
	}
	n, err := syscall.Write(fd, b)
	if err != nil {
		err = fmt.Errorf("Write error %d: %s", of, err)
		return 0, err
	}
	return n, nil
}

func eraseQSPI(of uint32, sz uint32) error {
	ei.length = mi.erasesize
	end := of + sz
	for ei.start = of; ei.start < end; ei.start += ei.length {
		fmt.Println("Erasing Block...", ei.start)
		_, _, e := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd),
			uintptr(MEMERASE), uintptr(unsafe.Pointer(ei)))
		if e != 0 {
			err := fmt.Errorf("Erase error %x: %s", ei.start, e)
			return err
		}
	}
	return nil
}

func UpdateEnv() (err error) {
	b, err := getPerFile()
	s := strings.Split(string(b), "\x00")
	ip := s[0]
	if len(string(ip)) > 500 {
		err = fmt.Errorf("no 'ip=' in per blk, skipping env update")
		return err
	}
	e, bootargs, err := GetEnv()
	if err != nil {
		return err
	}
	if !strings.Contains(e[bootargs], "ip=") {
		err = fmt.Errorf("no 'ip=' in env blk, skipping env update")
		return err
	}
	n := strings.SplitAfter(e[bootargs], "ip=")
	if n[1] == string(ip) {
		err = fmt.Errorf("no ip change, skipping env update")
		return err
	}
	e[bootargs] = n[0] + string(ip)
	err = PutEnv(e)
	if err != nil {
		return err
	}
	return nil
}

func GetEnv() (env []string, bootargs int, err error) {
	b, err := readBlk("env")
	if err != nil {
		return nil, 0, err
	}
	e := strings.Split(string(b[ENVCRC:ENVSIZE]), "\x00")
	var end int
	for j, n := range e {
		if strings.Contains(n, "bootargs") {
			bootargs = j
		}
		if len(n) == 0 {
			end = j
			break
		}
	}
	return e[:end], bootargs, nil
}

func PutEnv(e []string) (err error) {
	ee := strings.Join(e, "\x00")
	b := make([]byte, ENVSIZE, ENVSIZE)
	b = []byte(ee)
	for i := len(b); i < ENVSIZE; i++ {
		b = append(b, 0)
	}

	x := crc32.ChecksumIEEE(b[0 : ENVSIZE-ENVCRC])
	y := make([]byte, 4)
	binary.LittleEndian.PutUint32(y, x)
	b = append(y[0:4], b[0:ENVSIZE-ENVCRC]...)

	err = writeBlk("env", b[0:ENVSIZE])
	if err != nil {
		return err
	}
	return nil
}
