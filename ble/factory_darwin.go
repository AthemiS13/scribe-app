//go:build darwin

package ble

import "fmt"

// createPlatformBackend creates the macOS-specific BLE backend
func createPlatformBackend() (BLEBackend, error) {
	fmt.Println("Creating macOS BLE backend (optimized for macOS)")
	return NewMacOSBLEBackend(), nil
}
