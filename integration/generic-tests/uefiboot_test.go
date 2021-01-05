// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!race

package integration

import (
	"os"
	"testing"
	"time"

	expect "github.com/google/goexpect"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestUefiboot tests uefiboot commmands to boot to uefishell.
func TestUefiBoot(t *testing.T) {
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	img := "testdata/uefiboot/UEFIPAYLOAD.fd"
	if _, err := os.Stat(img); err != nil && os.IsNotExist(err) {
		t.Fatal("UEFI payload  image is not found.")
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: []string{"uefiboot /dev/sda"},
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				qemu.IDEBlockDevice{File: img},
				qemu.ArbitraryArgs{"-machine", "q35"},
				qemu.ArbitraryArgs{"-m", "2048"},
			},
		},
	})
	defer cleanup()

	if out, err := q.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: "PROGRESS CODE: V02070003"},
		// Last code before booting to UEFI Shell
		&expect.BExp{R: "PROGRESS CODE: V03058001"},
	}, 50*time.Second); err != nil {
		t.Fatalf("VM output did not match expectations: %v", out)
	}
}
