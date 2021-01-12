package smbios

import (
	"testing"
)

func TestGetSMBIOSEFI(t *testing.T) {
	systabPath = "testdata/systab"
	base, _, err := GetSMBIOSBaseEFI()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x12345678

	if base != want {
		t.Errorf("GetSMBIOSEFI() get 0x%x, want 0x%x", base, want)
	}
}
