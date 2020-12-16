// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/u-root/u-root/pkg/boot/uefi"
)

var (
	dryRun    = flag.Bool("dryrun", false, "Do not actually kexec into the boot config")
	debug     = flag.Bool("d", false, "Print debug output")
	imageBase = flag.Uint64("i", 0x800000, "Where to load payload image")
)

var v = func(string, ...interface{}) {}

func main() {
	flag.Parse()
	v = log.Printf
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: uefiboot <payload>")
	}
	fv, err := uefi.New(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	fv.ImageBase = uintptr(*imageBase)

	if err := fv.Load(*debug); err != nil {
		log.Fatal(err)
	}
}
