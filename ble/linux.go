//go:build linux

package ble

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	serviceUUID            = "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
	dataCharacteristicUUID = "1c95d5e3-d8f7-413a-bf3d-7a2e5d7be87e"
	infoCharacteristicUUID = "beb5483e-36e1-4688-b7f5-ea07361b26a8"
	dataTransferVersion    = uint8(1)
)

// LinuxBLEBackend implements BLEBackend for Linux using TinyGo bluetooth
type LinuxBLEBackend struct {
	adapter         *bluetooth.Adapter
	connectedDevice *bluetooth.Device
	dataChar        *bluetooth.DeviceCharacteristic
	infoChar        *bluetooth.DeviceCharacteristic
}

// NewLinuxBLEBackend creates a new Linux BLE backend
func NewLinuxBLEBackend() *LinuxBLEBackend {
	return &LinuxBLEBackend{}
}

// Initialize initializes the BLE adapter
func (l *LinuxBLEBackend) Initialize() error {
	l.adapter = bluetooth.DefaultAdapter
	if l.adapter == nil {
		return ErrAdapterNotFound
	}

	err := l.adapter.Enable()
	if err != nil {
		return err
	}

	return nil
}

// ScanForDevice scans for a device with the given name
func (l *LinuxBLEBackend) ScanForDevice(deviceName string, timeoutSeconds int) (*DeviceInfo, error) {
	if l.adapter == nil {
		return nil, ErrAdapterNotFound
	}

	ch := make(chan bluetooth.ScanResult, 1)
	var scanErr error

	go func() {
		err := l.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.LocalName() == deviceName {
				adapter.StopScan()
				ch <- device
			}
		})
		if err != nil {
			scanErr = err
		}
	}()

	select {
	case device := <-ch:
		return &DeviceInfo{
			Name:    device.LocalName(),
			Address: device.Address.String(),
		}, nil
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		l.adapter.StopScan()
		if scanErr != nil {
			return nil, scanErr
		}
		return nil, ErrDeviceNotFound
	}
}

// Connect connects to a device using its address
func (l *LinuxBLEBackend) Connect(deviceAddress string) error {
	if l.adapter == nil {
		return ErrAdapterNotFound
	}

	macAddress := bluetooth.Address{}
	macAddress.Set(deviceAddress)

	device, err := l.adapter.Connect(macAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return ErrConnectionFailed
	}

	// Discover services
	serviceID, err := bluetooth.ParseUUID(serviceUUID)
	if err != nil {
		return err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{serviceID})
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return errors.New("service not found")
	}

	service := services[0]

	// Discover characteristics
	dataCharUUID, err := bluetooth.ParseUUID(dataCharacteristicUUID)
	if err != nil {
		return err
	}

	infoCharUUID, err := bluetooth.ParseUUID(infoCharacteristicUUID)
	if err != nil {
		return err
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{dataCharUUID, infoCharUUID})
	if err != nil {
		return err
	}

	if len(chars) != 2 {
		return errors.New("characteristics not found")
	}

	// Assign characteristics
	if chars[0].UUID().String() == dataCharacteristicUUID {
		l.dataChar = &chars[0]
		l.infoChar = &chars[1]
	} else {
		l.dataChar = &chars[1]
		l.infoChar = &chars[0]
	}

	l.connectedDevice = &device
	return nil
}

// Disconnect disconnects from the currently connected device
func (l *LinuxBLEBackend) Disconnect() error {
	if l.connectedDevice == nil {
		return ErrNotConnected
	}

	err := l.connectedDevice.Disconnect()

	// Clear references regardless of error
	l.connectedDevice = nil
	l.dataChar = nil
	l.infoChar = nil

	return err
}

// IsConnected returns true if currently connected to a device
func (l *LinuxBLEBackend) IsConnected() bool {
	return l.connectedDevice != nil && l.dataChar != nil && l.infoChar != nil
}

// SendData sends data to the connected device
func (l *LinuxBLEBackend) SendData(data string) error {
	if !l.IsConnected() {
		return ErrNotConnected
	}

	// Send info first
	_, err := l.infoChar.WriteWithoutResponse(l.infoToBytes(dataTransferVersion, uint16(len(data)), 1))
	if err != nil {
		return ErrSendFailed
	}

	// Allow SCRIBE firmware to process info packet and allocate buffer
	time.Sleep(150 * time.Millisecond)

	// Send data in chunks
	chunkSize := 20
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		_, err = l.dataChar.WriteWithoutResponse([]byte(data[i:end]))
		if err != nil {
			return ErrSendFailed
		}

		time.Sleep(5 * time.Millisecond)
	}

	// Wait for SCRIBE to finish processing complete transmission
	// This ensures to_read = -1 before any subsequent SendData() calls
	processingTime := time.Duration(len(data)/10+400) * time.Millisecond // Base 400ms + processing time
	if processingTime < 600*time.Millisecond {
		processingTime = 600 * time.Millisecond // Minimum wait time
	}

	fmt.Printf("Waiting %v for SCRIBE to complete processing...\n", processingTime)
	time.Sleep(processingTime)

	return nil
}

// Close closes the BLE backend and releases resources
func (l *LinuxBLEBackend) Close() error {
	if l.IsConnected() {
		return l.Disconnect()
	}
	return nil
}

/*
data transfer info
	- total size of info: 4 bytes
	- 1st byte = version
	- 2nd byte = font size
	- 3-4th byte = data size
*/

// infoToBytes converts transfer info to bytes
func (l *LinuxBLEBackend) infoToBytes(v uint8, size uint16, font uint8) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(v)
	buf[1] = byte(font)
	binary.LittleEndian.PutUint16(buf[2:], size)
	return buf
}
