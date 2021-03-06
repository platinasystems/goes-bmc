// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import "github.com/platinasystems/goes-bmc/cmd/fantrayd"

func fantraydInit() {
	fantrayd.Vdev.Bus = 14
	fantrayd.Vdev.Addr = 0x20

	fantrayd.VpageByKey = map[string]uint8{
		"fan_tray.1.status": 1,
		"fan_tray.2.status": 2,
		"fan_tray.3.status": 3,
		"fan_tray.4.status": 4,
	}

	fantrayd.WrRegDv["fantrayd"] = "fantrayd"
	fantrayd.WrRegFn["fantrayd.example"] = "example"
	fantrayd.WrRegRng["fantrayd.example"] = []string{"true", "false"}
}
