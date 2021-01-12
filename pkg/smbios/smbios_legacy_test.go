package smbios

import (
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/memio"
)

var tmpBuf = []byte{0, 0, 0, 0, 0, 0}

func mockMemioRead(base int64, uintn memio.UintN) error {
	dat, ok := uintn.(*memio.ByteSlice)
	if !ok {
		return fmt.Errorf("Not supported.")
	}
	bufLen := len(tmpBuf)
	for i := int64(0); i < dat.Size(); i++ {
		(*dat)[i] = tmpBuf[(base+i)%int64(bufLen)]
	}
	return nil
}

func TestGetSMBIOSLegacyNotFound(t *testing.T) {
	base, err := GetSMBIOSBaseLegacy()
	if err == nil {
		t.Errorf("GetSMBIOSLegacy() get 0x%x, want NotFound", base)
	}
}

func TestGetSMBIOSLegacySMBIOS(t *testing.T) {
	tmpBuf = []byte{0, '_', 'M', 'S', '_', 0, 0, '_', 'S', 'M', '_', 0, 0, 0, 0, 0}
	memioRead = mockMemioRead
	base, err := GetSMBIOSBaseLegacy()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0xf0007

	if base != want {
		t.Errorf("GetSMBIOSLegacy() get 0x%x, want 0x%x", base, want)
	}
}

func TestGetSMBIOSLegacySMBIOS3(t *testing.T) {
	tmpBuf = []byte{0, '_', 'M', 'S', '_', 0, 0, '_', 'S', 'M', '3', '_', 0, 0, 0, 0, 0}
	memioRead = mockMemioRead
	base, err := GetSMBIOSBaseLegacy()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0xf0009

	if base != want {
		t.Errorf("GetSMBIOSLegacy() get 0x%x, want 0x%x", base, want)
	}
}
