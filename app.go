package main

import (
	"context"
	"fmt"
	"scribe-app/communication"

	"tinygo.org/x/bluetooth"
)

// App struct
type App struct {
	ctx             context.Context
	adapter         *bluetooth.Adapter
	connectedDevice *bluetooth.Device
	dataChar        *bluetooth.DeviceCharacteristic
	infoChar        *bluetooth.DeviceCharacteristic
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.adapter = bluetooth.DefaultAdapter
	err := a.adapter.Enable()
	if err != nil {
		fmt.Println(err)
		a.adapter = nil
	}
	a.dataChar = nil
	a.infoChar = nil
	a.connectedDevice = nil
}

func (a *App) Connect() bool {
	device, dataChar, infoChar, err := communication.ConnectDevice(a.adapter)

	if err != nil {
		fmt.Println(err)
		return false
	}

	a.connectedDevice = &device
	a.dataChar = &dataChar
	a.infoChar = &infoChar

	fmt.Println("SCRIBE connected")

	return true
}

func (a *App) Disconnect() bool {
	err := communication.DisconnectDevice(*a.connectedDevice)

	if err != nil {
		fmt.Println(err)
		return false
	}

	a.connectedDevice = nil
	a.dataChar = nil
	a.infoChar = nil

	fmt.Println("SCRIBE disconnected")

	return true
}

func (a *App) SendData(data string) bool {
	err := communication.SendData(*a.dataChar, *a.infoChar, data)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
