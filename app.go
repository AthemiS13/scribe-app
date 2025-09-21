package main

import (
	"context"
	"fmt"
	"scribe-app/ble"
)

// App struct
type App struct {
	ctx        context.Context
	bleBackend ble.BLEBackend
	factory    *ble.Factory
	deviceInfo *ble.DeviceInfo
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		factory: ble.NewFactory(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Display platform information
	currentPlatform, info := a.factory.GetPlatformInfo()
	fmt.Printf("Platform: %s\n", currentPlatform.String())
	fmt.Printf("BLE Info: %s\n", info)

	// Create platform-specific BLE backend
	backend, err := a.factory.CreateBackend()
	if err != nil {
		fmt.Printf("Failed to create BLE backend: %v\n", err)
		return
	}

	a.bleBackend = backend

	// Initialize the BLE backend
	err = a.bleBackend.Initialize()
	if err != nil {
		fmt.Printf("Failed to initialize BLE backend: %v\n", err)
		a.bleBackend = nil
	} else {
		fmt.Println("BLE backend initialized successfully")
	}
}

func (a *App) Connect() bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	// First scan for the device
	deviceInfo, err := a.bleBackend.ScanForDevice("SCRIBE", 10)
	if err != nil {
		fmt.Printf("Failed to find device: %v\n", err)
		return false
	}

	a.deviceInfo = deviceInfo
	fmt.Printf("Found device: %s at %s\n", deviceInfo.Name, deviceInfo.Address)

	// Connect to the device
	err = a.bleBackend.Connect(deviceInfo.Address)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return false
	}

	fmt.Println("SCRIBE connected")
	return true
}

func (a *App) Disconnect() bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	err := a.bleBackend.Disconnect()
	if err != nil {
		fmt.Printf("Failed to disconnect: %v\n", err)
		return false
	}

	a.deviceInfo = nil
	fmt.Println("SCRIBE disconnected")
	return true
}

func (a *App) SendData(data string) bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	if !a.bleBackend.IsConnected() {
		fmt.Println("Not connected to device")
		return false
	}

	fmt.Printf("=== SENDING DATA WHILE STAYING CONNECTED ===\n")
	fmt.Printf("Connection state: %v\n", a.bleBackend.IsConnected())

	err := a.bleBackend.SendData(data)
	if err != nil {
		fmt.Printf("Failed to send data: %v\n", err)
		return false
	}

	fmt.Printf("Data sent successfully: %s\n", data)
	fmt.Printf("=== TRANSMISSION COMPLETE - STAYING CONNECTED ===\n")
	return true
}

// GetPlatformInfo returns information about the current platform
func (a *App) GetPlatformInfo() string {
	if a.factory == nil {
		return "Factory not initialized"
	}

	currentPlatform, info := a.factory.GetPlatformInfo()
	return fmt.Sprintf("Platform: %s - %s", currentPlatform.String(), info)
}

// IsConnected returns the connection status
func (a *App) IsConnected() bool {
	if a.bleBackend == nil {
		return false
	}
	return a.bleBackend.IsConnected()
}

// ForceDisconnect explicitly disconnects (call this when completely done)
func (a *App) ForceDisconnect() bool {
	fmt.Println("=== FORCE DISCONNECT REQUESTED ===")
	return a.Disconnect()
}
