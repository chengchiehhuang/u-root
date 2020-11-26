// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/fit"
)

var (
	dryRun    = flag.Bool("dryrun", false, "Do not actually kexec into the boot config")
	debug     = flag.Bool("d", false, "Print debug output")
	cmdline = flag.String("c", "earlyprintk=ttyS0,115200,keep console=ttyS0", "command line")
)

var v = func(string, ...interface{}) {}

func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: fitboot fitimage")
	}
	f, err := fit.New(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	v("Loaded fitimage: %s", f)
	f.Cmdline = *cmdline
	if err := f.LoadFITImage(*debug); err != nil {
		log.Fatal(err)
	}
	if *dryRun {
		v("Not trying to boot since this is a dry run")
		os.Exit(0)
	}
	if err := boot.Execute(); err != nil {
		log.Fatal(err)
	}
}
