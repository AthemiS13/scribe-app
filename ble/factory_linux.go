//go:build linux

package ble

import "fmt"

// createPlatformBackend creates the Linux-specific BLE backend
func createPlatformBackend() (BLEBackend, error) {
	fmt.Println("Creating Linux BLE backend (using TinyGo bluetooth)")
	return NewLinuxBLEBackend(), nil
}
