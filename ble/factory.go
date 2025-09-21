package ble

import (
	"scribe-app/platform"
)

// Factory creates the appropriate BLE backend based on the current platform
type Factory struct{}

// NewFactory creates a new BLE factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateBackend creates and returns the appropriate BLE backend for the current platform
func (f *Factory) CreateBackend() (BLEBackend, error) {
	return createPlatformBackend()
}

// GetPlatformInfo returns information about the current platform and its BLE support
func (f *Factory) GetPlatformInfo() (platform.Platform, string) {
	currentPlatform := platform.DetectPlatform()

	var info string
	switch currentPlatform {
	case platform.Linux:
		info = "Linux platform using TinyGo bluetooth library for BLE communication"
	case platform.MacOS:
		info = "macOS platform with optimized BLE implementation for better reliability"
	case platform.Windows:
		info = "Windows platform with Windows BLE API integration for improved performance"
	default:
		info = "Unsupported platform"
	}

	return currentPlatform, info
}
