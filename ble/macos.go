//go:build darwin

package ble

import (
	"encoding/binary"
	"fmt"
	"os/exec"
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

// MacOSBLEBackend implements BLEBackend for macOS using TinyGo bluetooth with macOS optimizations
type MacOSBLEBackend struct {
	adapter          *bluetooth.Adapter
	connectedDevice  *bluetooth.Device
	dataChar         *bluetooth.DeviceCharacteristic
	infoChar         *bluetooth.DeviceCharacteristic
	deviceInfo       *DeviceInfo
	mutex            sync.RWMutex
	lastConnectTime  time.Time
	consecutiveSends int
}

// NewMacOSBLEBackend creates a new macOS BLE backend
func NewMacOSBLEBackend() *MacOSBLEBackend {
	return &MacOSBLEBackend{}
}

// Initialize initializes the BLE adapter
func (m *MacOSBLEBackend) Initialize() error {
	// Check if blueutil is available (helpful for macOS BLE management)
	_, err := exec.LookPath("blueutil")
	if err != nil {
		fmt.Println("Info: blueutil not found. Consider installing with: brew install blueutil")
	}

	m.adapter = bluetooth.DefaultAdapter
	if m.adapter == nil {
		return ErrAdapterNotFound
	}

	err = m.adapter.Enable()
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

	go func() {
		err := m.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.LocalName() == "SCRIBE" && device.Address.String() == deviceAddress {
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
		// Found the device
	case <-time.After(10 * time.Second):
		m.adapter.StopScan()
		if scanErr != nil {
			return scanErr
		}
		return ErrDeviceNotFound
	}

	// Connect with macOS-optimized parameters
	device, err := m.adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(10 * time.Second),
	})
	if err != nil {
		return fmt.Errorf("macOS connection failed: %v", err)
	}

	// Discover SCRIBE service
	serviceID, err := bluetooth.ParseUUID(scribeServiceUUID)
	if err != nil {
		return err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{serviceID})
	if err != nil {
		return fmt.Errorf("service discovery failed on macOS: %v", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("SCRIBE service not found")
	}

	service := services[0]

	// Discover SCRIBE characteristics
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
		return fmt.Errorf("characteristic discovery failed on macOS: %v", err)
	}

	if len(chars) != 2 {
		return fmt.Errorf("SCRIBE characteristics not found (found %d, expected 2)", len(chars))
	}

	// Assign characteristics based on UUID (case-insensitive comparison)
	for _, char := range chars {
		charUUID := strings.ToUpper(char.UUID().String())
		if charUUID == strings.ToUpper(scribeDataCharacteristicUUID) {
			m.dataChar = &char
			fmt.Printf("Found data characteristic: %s\n", char.UUID().String())
		} else if charUUID == strings.ToUpper(scribeInfoCharacteristicUUID) {
			m.infoChar = &char
			fmt.Printf("Found info characteristic: %s\n", char.UUID().String())
		}
	}

	if m.dataChar == nil || m.infoChar == nil {
		return fmt.Errorf("failed to identify SCRIBE characteristics")
	}

	m.mutex.Lock()
	m.connectedDevice = &device
	m.lastConnectTime = time.Now()
	m.consecutiveSends = 0
	m.mutex.Unlock()

	fmt.Println("Successfully connected to SCRIBE on macOS")
	return nil
}

// Disconnect disconnects from the currently connected device
func (m *MacOSBLEBackend) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.connectedDevice == nil {
		return ErrNotConnected
	}

	fmt.Println("Disconnecting from SCRIBE on macOS...")

	err := m.connectedDevice.Disconnect()

	// Clear references regardless of error
	m.connectedDevice = nil
	m.dataChar = nil
	m.infoChar = nil
	m.deviceInfo = nil
	m.consecutiveSends = 0

	return err
}

// IsConnected returns true if currently connected to a device
func (m *MacOSBLEBackend) IsConnected() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.connectedDevice != nil && m.dataChar != nil && m.infoChar != nil
}

// SendData sends data to the connected device using SCRIBE protocol with macOS optimizations
func (m *MacOSBLEBackend) SendData(data string) error {
	if !m.IsConnected() {
		return ErrNotConnected
	}

	m.mutex.Lock()
	m.consecutiveSends++
	sendCount := m.consecutiveSends
	timeSinceConnect := time.Since(m.lastConnectTime)
	m.mutex.Unlock()

	fmt.Printf("=== SCRIBE TRANSMISSION #%d (Connected for: %v) ===\n", sendCount, timeSinceConnect.Round(time.Second))
	fmt.Printf("Sending data to SCRIBE on macOS (length: %d bytes)...\n", len(data))

	// Send info packet first (SCRIBE protocol)
	infoData := m.createInfoPacket(uint16(len(data)))
	_, err := m.infoChar.WriteWithoutResponse(infoData)
	if err != nil {
		return fmt.Errorf("failed to send info packet: %v", err)
	}

	fmt.Printf("Sent info packet: [version=%d, font=%d, size=%d]\n", infoData[0], infoData[1], binary.LittleEndian.Uint16(infoData[2:]))

	// Critical delay: Allow SCRIBE firmware to process info packet and allocate buffer
	// The firmware needs time to malloc() the data buffer and set up state
	time.Sleep(200 * time.Millisecond)

	// Send data in optimized chunks for macOS
	chunkSize := 18 // Optimal for macOS BLE stack
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := []byte(data[i:end])
		_, err = m.dataChar.WriteWithoutResponse(chunk)
		if err != nil {
			return fmt.Errorf("failed to send data chunk: %v", err)
		}

		fmt.Printf("Sent chunk %d-%d (%d bytes)\n", i, end-1, len(chunk))

		// macOS-optimized delay between chunks
		time.Sleep(10 * time.Millisecond)
	}

	// CRITICAL: Wait for SCRIBE to finish processing complete transmission
	// This ensures to_read = -1 before any subsequent SendData() calls
	// Time needed for: data processing + page creation + display update + flash storage

	// AGGRESSIVE TIMING: Give ESP32 enough time for memory management with multiple pages
	// PROBLEM: ESP32 heap fragmentation when creating 4+ pages with malloc()
	// SOLUTION: Extended delays allow memory defragmentation and proper allocation

	baseTime := 800                  // Increased base time for memory operations
	dataComplexity := len(data) * 12 // Much more time for complex page parsing
	processingTime := time.Duration(baseTime+dataComplexity) * time.Millisecond

	// Extended limits for unlimited pages
	if len(data) <= 20 {
		// Short data: still allow reasonable time
		if processingTime > 1200*time.Millisecond {
			processingTime = 1200 * time.Millisecond
		}
	} else {
		// Longer data: give ESP32 plenty of time for memory management
		if processingTime > 3000*time.Millisecond {
			processingTime = 3000 * time.Millisecond
		}
	}

	// Ensure minimum time for memory operations
	if processingTime < 800*time.Millisecond {
		processingTime = 800 * time.Millisecond
	}

	fmt.Printf("Processing delay: %v (extended timing for unlimited pages)\n", processingTime)
	time.Sleep(processingTime)

	// Warn about potential firmware state issues after multiple consecutive sends
	if sendCount >= 3 {
		fmt.Printf("⚠️  WARNING: This is transmission #%d without disconnecting.\n", sendCount)
		fmt.Println("⚠️  SCRIBE firmware may be in an unstable state. Consider disconnecting/reconnecting.")
	}

	fmt.Printf("Data sent successfully to SCRIBE via macOS BLE (Total transmissions: %d)\n", sendCount)
	return nil
}

/*
data transfer info
	- total size of info: 4 bytes
	- 1st byte = version
	- 2nd byte = font size
	- 3-4th byte = data size
*/

// createInfoPacket creates the SCRIBE info packet (4 bytes: version, font, data_size_low, data_size_high)
func (m *MacOSBLEBackend) createInfoPacket(dataSize uint16) []byte {
	buf := make([]byte, 4)
	buf[0] = scribeDataTransferVersion               // version = 1
	buf[1] = 1                                       // font size = 1
	binary.LittleEndian.PutUint16(buf[2:], dataSize) // little-endian data size
	return buf
}

// Close closes the BLE backend and releases resources
func (m *MacOSBLEBackend) Close() error {
	return m.Disconnect()
}

// getMacOSBluetoothStatus checks the Bluetooth status on macOS
func (m *MacOSBLEBackend) getMacOSBluetoothStatus() (bool, error) {
	cmd := exec.Command("system_profiler", "SPBluetoothDataType")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "State: On"), nil
}
