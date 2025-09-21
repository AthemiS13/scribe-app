package ble

import (
	"errors"
)

// Common errors
var (
	ErrDeviceNotFound   = errors.New("device not found")
	ErrConnectionFailed = errors.New("connection failed")
	ErrNotConnected     = errors.New("not connected")
	ErrAdapterNotFound  = errors.New("no ble adapter found")
	ErrTimeout          = errors.New("timeout exceeded")
	ErrSendFailed       = errors.New("send failed")
)

// DeviceInfo represents information about a BLE device
type DeviceInfo struct {
	Name    string
	Address string
}

// BLEBackend defines the interface that all platform-specific BLE implementations must implement
type BLEBackend interface {
	// Initialize initializes the BLE adapter
	Initialize() error

	// ScanForDevice scans for a device with the given name and returns its info
	ScanForDevice(deviceName string, timeoutSeconds int) (*DeviceInfo, error)

	// Connect connects to a device using its address
	Connect(deviceAddress string) error

	// Disconnect disconnects from the currently connected device
	Disconnect() error

	// IsConnected returns true if currently connected to a device
	IsConnected() bool

	// SendData sends data to the connected device
	SendData(data string) error

	// Close closes the BLE backend and releases resources
	Close() error
}
