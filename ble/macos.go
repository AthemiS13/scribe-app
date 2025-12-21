//go:build darwin

package ble

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	// SCRIBE firmware BLE UUIDs from main.cpp
	scribeServiceUUID            = "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
	scribeInfoCharacteristicUUID = "beb5483e-36e1-4688-b7f5-ea07361b26a8"
	scribeDataCharacteristicUUID = "1c95d5e3-d8f7-413a-bf3d-7a2e5d7be87e"
	scribeDataTransferVersion    = uint8(1)
)

// STATUS_CODES from firmware
const (
	STATUS_READY = 0x01
	STATUS_DONE  = 0x02
)

// MacOSBLEBackend implements BLEBackend for macOS using TinyGo bluetooth with notifications
type MacOSBLEBackend struct {
	adapter         *bluetooth.Adapter
	connectedDevice *bluetooth.Device
	dataChar        *bluetooth.DeviceCharacteristic
	infoChar        *bluetooth.DeviceCharacteristic
	deviceInfo      *DeviceInfo

	statusChan chan uint8
	mutex      sync.RWMutex
}

// NewMacOSBLEBackend creates a new macOS BLE backend
func NewMacOSBLEBackend() *MacOSBLEBackend {
	return &MacOSBLEBackend{
		statusChan: make(chan uint8, 10),
	}
}

// Initialize initializes the BLE adapter
func (m *MacOSBLEBackend) Initialize() error {
	m.adapter = bluetooth.DefaultAdapter
	if m.adapter == nil {
		return ErrAdapterNotFound
	}

	err := m.adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable BLE adapter on macOS: %v", err)
	}

	fmt.Println("macOS BLE adapter initialized successfully")
	return nil
}

// ScanForDevice scans for a device with the given name
func (m *MacOSBLEBackend) ScanForDevice(deviceName string, timeoutSeconds int) (*DeviceInfo, error) {
	if m.adapter == nil {
		return nil, ErrAdapterNotFound
	}

	fmt.Printf("Scanning for '%s' on macOS (timeout: %d seconds)...\n", deviceName, timeoutSeconds)

	ch := make(chan bluetooth.ScanResult, 1)
	var scanErr error

	go func() {
		err := m.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.LocalName() == deviceName {
				fmt.Printf("Found device: %s at %s\n", device.LocalName(), device.Address.String())
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
		deviceInfo := &DeviceInfo{
			Name:    device.LocalName(),
			Address: device.Address.String(),
		}
		m.deviceInfo = deviceInfo
		return deviceInfo, nil
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		m.adapter.StopScan()
		if scanErr != nil {
			return nil, fmt.Errorf("scan error on macOS: %v", scanErr)
		}
		return nil, ErrDeviceNotFound
	}
}

// Connect connects to a device using its address
func (m *MacOSBLEBackend) Connect(deviceAddress string) error {
	if m.adapter == nil {
		return ErrAdapterNotFound
	}

	if m.deviceInfo == nil || m.deviceInfo.Address != deviceAddress {
		return ErrDeviceNotFound
	}

	fmt.Printf("Connecting to SCRIBE at %s on macOS...\n", deviceAddress)

	// Re-scan to get fresh ScanResult for connection
	ch := make(chan bluetooth.ScanResult, 1)
	var scanErr error

	// Short scan to energize the adapter and find the peripheral
	go func() {
		err := m.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.Address.String() == deviceAddress {
				adapter.StopScan()
				ch <- device
			}
		})
		if err != nil {
			scanErr = err
		}
	}()

	var scanResult bluetooth.ScanResult
	select {
	case scanResult = <-ch:
		// Found
	case <-time.After(5 * time.Second):
		m.adapter.StopScan()
		if scanErr != nil {
			return scanErr
		}
		return ErrDeviceNotFound
	}

	device, err := m.adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(10 * time.Second),
	})
	if err != nil {
		return fmt.Errorf("macOS connection failed: %v", err)
	}

	// Discover Services
	serviceID, err := bluetooth.ParseUUID(scribeServiceUUID)
	if err != nil {
		return err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{serviceID})
	if err != nil {
		return fmt.Errorf("service discovery failed: %v", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("SCRIBE service not found")
	}

	service := services[0]

	// Discover Characteristics
	dataCharUUID, err := bluetooth.ParseUUID(scribeDataCharacteristicUUID)
	if err != nil {
		return err
	}
	infoCharUUID, err := bluetooth.ParseUUID(scribeInfoCharacteristicUUID)
	if err != nil {
		return err
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{dataCharUUID, infoCharUUID})
	if err != nil {
		return fmt.Errorf("characteristic discovery failed: %v", err)
	}

	if len(chars) != 2 {
		return fmt.Errorf("SCRIBE characteristics not found (found %d, expected 2)", len(chars))
	}

	m.mutex.Lock()
	m.connectedDevice = &device

	// Assign chars
	for _, char := range chars {
		if strings.EqualFold(char.UUID().String(), scribeDataCharacteristicUUID) {
			c := char
			m.dataChar = &c
		} else if strings.EqualFold(char.UUID().String(), scribeInfoCharacteristicUUID) {
			c := char
			m.infoChar = &c
		}
	}
	m.mutex.Unlock()

	if m.dataChar == nil || m.infoChar == nil {
		return fmt.Errorf("failed to identify SCRIBE characteristics")
	}

	// Enable Notifications on InfoChar for flow control
	fmt.Println("Enabling notifications on Info Characteristic...")
	err = m.infoChar.EnableNotifications(func(buf []byte) {
		if len(buf) > 0 {
			msg := buf[0]
			fmt.Printf("🔔 NOTIFICATION RECEIVED: 0x%02x\n", msg)
			// Non-blocking send to channel
			select {
			case m.statusChan <- msg:
			default:
				fmt.Println("⚠️ Status channel full, dropping notification")
			}
		}
	})
	if err != nil {
		fmt.Printf("Warning: Failed to enable notifications: %v. Flow control may fail.\n", err)
	}

	fmt.Println("Successfully connected to SCRIBE on macOS with Notification support")
	return nil
}

// Disconnect disconnects from the currently connected device
func (m *MacOSBLEBackend) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.connectedDevice == nil {
		return ErrNotConnected
	}

	fmt.Println("Disconnecting from SCRIBE...")

	// Drain channel
Loop:
	for {
		select {
		case <-m.statusChan:
		default:
			break Loop
		}
	}

	err := m.connectedDevice.Disconnect()
	m.connectedDevice = nil
	m.dataChar = nil
	m.infoChar = nil
	m.deviceInfo = nil
	return err
}

// IsConnected returns true if currently connected
func (m *MacOSBLEBackend) IsConnected() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.connectedDevice != nil
}

// SendData sends data using robust flow control
func (m *MacOSBLEBackend) SendData(data string) error {
	if !m.IsConnected() {
		return ErrNotConnected
	}

	// clear channel before starting
Loop:
	for {
		select {
		case <-m.statusChan:
		default:
			break Loop
		}
	}

	dataLen := len(data)
	fmt.Printf("🔵 STARTING TRANSMISSION: %d bytes\n", dataLen)

	// 1. Send Info Packet
	infoData := m.createInfoPacket(uint16(dataLen))
	fmt.Printf("➡️ Sending Info Packet (expecting READY signal)...\n")
	// Use Write (confirmed) for the Info packet to ensure handshake starts correctly
	_, err := m.infoChar.Write(infoData)
	if err != nil {
		return fmt.Errorf("failed to send info packet: %v", err)
	}

	// 2. Wait for READY (0x01)
	// We might have old 0x02s floating around, so we should be specific
WaitReady:
	for {
		select {
		case status := <-m.statusChan:
			if status == STATUS_READY {
				fmt.Println("✅ RECEIVED READY SIGNAL (0x01)")
				break WaitReady
			} else {
				fmt.Printf("⚠️ Ignoring unexpected status while waiting for READY: 0x%02x\n", status)
			}
		case <-time.After(5 * time.Second):
			return fmt.Errorf("timeout waiting for READY signal from Scribe")
		}
	}

	// 3. Send Data Chunks
	// CRITICAL: WriteWithoutResponse + 30ms Pacing to avoid Duplicates (Retries) and Drops (Overflow)
	chunkSize := 18 // Reduced Chunk Size
	totalChunks := (dataLen + chunkSize - 1) / chunkSize

	fmt.Printf("➡️ Sending %d chunks...\n", totalChunks)

	for i := 0; i < dataLen; i += chunkSize {
		end := i + chunkSize
		if end > dataLen {
			end = dataLen
		}
		chunk := []byte(data[i:end])

		// WriteWithoutResponse is faster and doesn't auto-retry (preventing duplicates)
		_, err := m.dataChar.WriteWithoutResponse(chunk)
		if err != nil {
			return fmt.Errorf("failed to send chunk %d: %v", i/chunkSize, err)
		}

		// 30ms delay allows firmware loop() to process previous chunk reliably
		time.Sleep(30 * time.Millisecond)
	}

	fmt.Println("➡️ All chunks sent. Waiting for confirmation...")

	// 4. Wait for DONE (0x02)
WaitDone:
	for {
		select {
		case status := <-m.statusChan:
			if status == STATUS_DONE {
				fmt.Println("✅ RECEIVED DONE SIGNAL (0x02)")
				break WaitDone
			} else if status == STATUS_READY {
				fmt.Println("⚠️ Received late READY signal while waiting for DONE (ignoring)")
			} else {
				fmt.Printf("⚠️ Unexpected status: 0x%02x\n", status)
			}
		case <-time.After(5 * time.Second):
			return fmt.Errorf("timeout waiting for DONE signal from Scribe")
		}
	}

	fmt.Println("🚀 TRANSMISSION SUCCESSFUL")
	return nil
}

// createInfoPacket creates the SCRIBE info packet
func (m *MacOSBLEBackend) createInfoPacket(dataSize uint16) []byte {
	buf := make([]byte, 4)
	buf[0] = scribeDataTransferVersion
	buf[1] = 1 // font size
	binary.LittleEndian.PutUint16(buf[2:], dataSize)
	return buf
}

// Close closes the BLE backend
func (m *MacOSBLEBackend) Close() error {
	return m.Disconnect()
}
