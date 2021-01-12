// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
)

const fvEntryImageOffset int64 = 0xA0

var bootExecute = boot.Execute
var kexecLoad = kexec.Load
var kexecParseMemoryMap = kexec.ParseMemoryMap
var getRSDP = acpi.GetRSDP

func getSMBIOSBase() (int64, int64, error) {
	base, size, err := smbios.GetSMBIOSBaseEFI()
	if err != nil {
		base, size, err = smbios.GetSMBIOSBaseLegacy()
		if err != nil {
			return 0, 0, err
		}
	}
	return base, size, nil
}

// Serial port setting in Linuxboot.h
const (
	SerialPortTypeIO   = 1
	SerialPortTypeMMIO = 2
)

type SerialPortConfig struct {
	Type        uint32
	BaseAddr    uint32
	Baud        uint32
	RegWidth    uint32
	InputHertz  uint32
	UartPciAddr uint32
}

type payloadConfig struct {
	AcpiBase            uint64
	AcpiSize            uint64
	SmbiosBase          uint64
	SmbiosSize          uint64
	SerialConfig        SerialPortConfig
	NumMemoryMapEntries uint32
}

// FvImage is a structure for loading a firmware volume
type FvImage struct {
	name         string
	mem          kexec.Memory
	entryAddress uintptr
	ImageBase    uintptr
	SerialConfig SerialPortConfig
}

// PEImage is a structure for loading entry PE image from a firmware volume
type PEImage struct {
	file *os.File
}

// ReadAt of PEImage will add fvEntryImageOffset to regular ReadAt's offset
func (fv *PEImage) ReadAt(p []byte, off int64) (n int, err error) {
	// Entry PE image for UEFI Payload is 0xA0 offset from FV head.
	return fv.file.ReadAt(p, off+fvEntryImageOffset)
}

func checkFvAndGetEntryPoint(name string) (uintptr, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	pf, err := pe.NewFile(&PEImage{file: f})
	if err != nil {
		return 0, err
	}
	return uintptr(fvEntryImageOffset) + uintptr(pf.OptionalHeader.(*pe.OptionalHeader64).AddressOfEntryPoint), nil
}

// New loads the file and return FvImage stucture if entry image is found
func New(n string) (*FvImage, error) {
	entry, err := checkFvAndGetEntryPoint(n)
	if err != nil {
		return nil, err
	}
	return &FvImage{name: n, mem: kexec.Memory{}, entryAddress: entry}, nil
}

// Reserved 64kb for passing params
const uefiPayloadConfigSize = 0x10000

// Load loads fimware volume payload and boot the the payload
func (fv *FvImage) Load(verbose bool) error {
	// Install payload
	dat, err := ioutil.ReadFile(fv.name)
	if err != nil {
		return err
	}
	fv.mem.Segments.Insert(kexec.NewSegment(dat, kexec.Range{Start: fv.ImageBase, Size: uint(len(dat))}))

	// Install payload config & its memory map: 64 kb below the image
	// We cannot use the memory above  the image base because it may be used by HOB
	var configAddr uintptr = fv.ImageBase - uintptr(uefiPayloadConfigSize)

	// Get MemoryMap
	mm, err := kexecParseMemoryMap()
	if err != nil {
		return err
	}

	// Get Acpi Basc (RSDP)
	rsdp, err := getRSDP()
	if err != nil {
		return err
	}

	smbiosBase, smbiosSize, err := getSMBIOSBase()
	if err != nil {
		return err
	}

	pc := payloadConfig{
		AcpiBase:            uint64(rsdp.RSDPAddr()),
		AcpiSize:            uint64(rsdp.Len()),
		SmbiosBase:          uint64(smbiosBase),
		SmbiosSize:          uint64(smbiosSize),
		SerialConfig:        fv.SerialConfig,
		NumMemoryMapEntries: uint32(len(mm)),
	}

	pcbuf := &bytes.Buffer{}
	if err := binary.Write(pcbuf, binary.LittleEndian, pc); err != nil {
		return err
	}

	if err := binary.Write(pcbuf, binary.LittleEndian, mm.AsPayloadParam()); err != nil {
		return err
	}

	if len(pcbuf.Bytes()) > uefiPayloadConfigSize {
		return fmt.Errorf("Config/Memmap size is greater than reserved size: %d bytes", len(pcbuf.Bytes()))
	}

	fv.mem.Segments.Insert(kexec.NewSegment(pcbuf.Bytes(), kexec.Range{Start: configAddr, Size: uint(len(pcbuf.Bytes()))}))

	if verbose {
		log.Printf("segments cmdline %v", fv.mem.Segments)
	}

	if err := kexecLoad(fv.ImageBase+uintptr(fv.entryAddress), fv.mem.Segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}

	if err := bootExecute(); err != nil {
		log.Fatal(err)
	}
	return nil
}
