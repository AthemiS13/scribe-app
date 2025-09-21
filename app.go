package main

import (
	"context"
	"fmt"
	"scribe-app/ble"
	"strings"
	"time"
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

// SendMultiplePages sends multiple data transmissions while staying connected
// This is more efficient and reliable than disconnecting/reconnecting
func (a *App) SendMultiplePages(pages []string) bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	if !a.bleBackend.IsConnected() {
		fmt.Println("Not connected to device")
		return false
	}

	fmt.Printf("Sending %d pages to SCRIBE...\n", len(pages))

	for i, page := range pages {
		fmt.Printf("Sending page %d/%d (length: %d bytes)\n", i+1, len(pages), len(page))

		err := a.bleBackend.SendData(page)
		if err != nil {
			fmt.Printf("Failed to send page %d: %v\n", i+1, err)
			return false
		}

		fmt.Printf("Page %d/%d sent successfully\n", i+1, len(pages))

		// Brief delay between pages to allow firmware processing without blocking button input
		if i < len(pages)-1 {
			fmt.Println("Brief pause between pages...")
			time.Sleep(200 * time.Millisecond) // Short enough to preserve button responsiveness
		}
	}

	fmt.Printf("All %d pages sent successfully!\n", len(pages))
	return true
}

// SendDataSmart automatically connects if needed and sends data
// Designed to be called multiple times without manual connection management
func (a *App) SendDataSmart(data string) bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	// Auto-connect if not connected
	if !a.bleBackend.IsConnected() {
		fmt.Println("Auto-connecting to SCRIBE...")
		if !a.Connect() {
			fmt.Println("Failed to auto-connect")
			return false
		}
	}

	fmt.Printf("=== SMART SEND: STAYING CONNECTED FOR MULTIPLE TRANSMISSIONS ===\n")

	err := a.bleBackend.SendData(data)
	if err != nil {
		fmt.Printf("Failed to send data: %v\n", err)
		return false
	}

	fmt.Printf("Data sent successfully (connection maintained): %s\n", data)
	return true
}

// SendAllPagesAsContinuousStream sends all pages as one continuous data stream
// Each page formatted to exactly 4 lines as firmware expects
func (a *App) SendAllPagesAsContinuousStream(pages []string) bool {
	if a.bleBackend == nil {
		fmt.Println("BLE backend not available")
		return false
	}

	if len(pages) == 0 {
		fmt.Println("No pages to send")
		return false
	}

	fmt.Printf("=== SENDING %d PAGES AS CONTINUOUS STREAM ===\n", len(pages))

	// Build the continuous data stream: each page must have exactly 4 lines
	var dataStream strings.Builder

	for i, page := range pages {
		if strings.TrimSpace(page) == "" {
			continue // Skip empty pages
		}

		// Add the page content
		dataStream.WriteString(page)

		// Count existing newlines and add more to make exactly 4 lines total
		newlineCount := strings.Count(page, "\n")
		// We need 3 more newlines after content to make 4 lines total
		// (content line + 3 empty lines = 4 lines per page)
		for newlineCount < 3 {
			dataStream.WriteString("\n")
			newlineCount++
		}

		fmt.Printf("Page %d: %q (4 lines)\n", i+1, page+strings.Repeat("\n", 3-strings.Count(page, "\n")))
	}

	finalData := dataStream.String()
	fmt.Printf("Final data stream: %q\n", finalData)
	fmt.Printf("Total length: %d bytes\n", len(finalData))

	// Connect once and send all data as continuous stream
	if !a.Connect() {
		fmt.Println("Failed to connect")
		return false
	}

	err := a.bleBackend.SendData(finalData)
	if err != nil {
		fmt.Printf("Failed to send continuous data stream: %v\n", err)
		a.Disconnect()
		return false
	}

	fmt.Printf("All %d pages sent as continuous stream!\n", len(pages))
	a.Disconnect()
	return true
} // ForceDisconnect explicitly disconnects (call this when completely done)
func (a *App) ForceDisconnect() bool {
	fmt.Println("=== FORCE DISCONNECT REQUESTED ===")
	return a.Disconnect()
}
