package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"
)

const envvar = `baudrate=115200
bootargs=console=ttymxc0,115200n8 ip=dhcp
bootcmd=run readmac;run sf_read_itb bootlinux_itb;run sf_read_legacy bootlinux_legacy
bootlinux_itb=bootm ${loadaddr}
bootlinux_legacy=bootz ${loadaddr} ${initrd_high} ${fdt_high}
bootdelay=3
ethact=FEC
ethprime=FEC
fdt_high=0x88000000
initrd_high=0x89000000
loadaddr=0x82000000
load_fdt_sf=sf read ${fdt_high} 0x00080000 ${sz_fdt}
load_initramfs_sf=sf read ${initrd_high} 0x00300000 ${sz_initrd}
load_kernel_sf=sf read ${loadaddr} 0x00100000 ${sz_kernel}
net_read_itb=${netbootmethod} ${loadaddr} ${serverpath}platina-mk1-bmc-itb.bin
netboot=run readmac net_read_itb bootlinux_itb
netbootmethod=dhcp
qspi0=mw 020e01b8 00000005; mw 20a8004 c7000000; mw 020a8000 4300ca05
qspi1=mw 020e01b8 00000005; mw 20a8004 c7000000; mw 020a8000 c300ca05
readmac=i2c read 55 0.2 200 80800000; setmac 80800000 24; saveenv
sf_read_itb=sf probe 0;sf read ${loadaddr} 0x00100000 ${sz_itb}
sf_read_legacy=sf probe 0;run load_kernel_sf load_initramfs_sf load_fdt_sf
stderr=serial
stdin=serial
stdout=serial
sz_fdt=f000
sz_initrd=300000
sz_kernel=200000
sz_itb=800000
wd=mw 020e01a0 00000005;mw 020e01a4 00000005;mw 020e01a8 00000005;mw 020e01b8 00000005;mw 020e01bc 00000005;mw 020a8000 0300ca05;mw 020a8004 07000000
`

const ubootEnvSize = 8192
const crcSize = 4
const endMarkerSize = 1

func makeUbootEnv() ([]byte, error) {
	if len(envvar) > (ubootEnvSize - (crcSize + endMarkerSize)) {
		return nil, fmt.Errorf("U-boot environment size %d is too large", len(envvar))
	}
	binenv := make([]byte, ubootEnvSize)
	copy(binenv[crc32.Size:], strings.Replace(envvar, "\n", "\x00", -1))
	crc := crc32.ChecksumIEEE(binenv[crc32.Size:])
	binary.LittleEndian.PutUint32(binenv[:crc32.Size], crc)
	return binenv, nil
}
