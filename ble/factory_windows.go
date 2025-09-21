//go:build windows

package ble

import "fmt"

// createPlatformBackend creates the Windows-specific BLE backend
func createPlatformBackend() (BLEBackend, error) {
	fmt.Println("Creating Windows BLE backend (optimized for Windows)")
	return NewWindowsBLEBackend(), nil
}
