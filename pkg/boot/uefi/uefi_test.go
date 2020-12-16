// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

type kexecLoadFunc func(entry uintptr, segments kexec.Segments, flags uint64) error

func TestLoadFvImage(t *testing.T) {
	fv, err := New("testdata/uefi.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = func() (*acpi.RSDP, error) {
		t.Log("mock acpi.GetRSDP()")
		return &acpi.RSDP{}, nil
	}
	defer func(old kexecLoadFunc) { kexecLoad = old }(kexecLoad)
	kexecLoad = func(entry uintptr, segments kexec.Segments, flags uint64) error {
		t.Log("mock kexec.Load()")
		return nil
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecParseMemoryMap = old }(kexecParseMemoryMap)
	kexecParseMemoryMap = func() (kexec.MemoryMap, error) {
		t.Log("mock kexec.ParseMemMap()")
		return kexec.MemoryMap{}, nil
	}
	defer func(old func() error) { bootExecute = old }(bootExecute)
	bootExecute = func() error {
		t.Log("mock boot.Execute()")
		return nil
	}

	if err = fv.Load(true); err != nil {
		t.Fatal(err)
	}
}
