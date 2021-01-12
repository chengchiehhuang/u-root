package smbios

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var systabPath = "/sys/firmware/efi/systab"

// GetSMBIOSBaseEFI finds the SMBIOS entry point address in the EFI System Table.
func GetSMBIOSBaseEFI() (int64, int64, error) {
	file, err := os.Open(systabPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	const (
		smbios3 = "SMBIOS3="
		smbios  = "SMBIOS="
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		start := ""
		size := int64(0)
		if strings.HasPrefix(line, smbios3) {
			start = strings.TrimPrefix(line, smbios3)
			size = 0x18
		}
		if strings.HasPrefix(line, smbios) {
			start = strings.TrimPrefix(line, smbios)
			size = 0x1f
		}
		if start == "" {
			continue
		}
		base, err := strconv.ParseInt(start, 0, 63)
		if err != nil {
			continue
		}
		return base, size, nil
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error while reading EFI systab: %v", err)
	}
	return 0, 0, fmt.Errorf("invalid /sys/firmware/efi/systab file")
}
