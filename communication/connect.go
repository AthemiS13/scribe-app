package communication

import (
	"errors"
	"time"

	"tinygo.org/x/bluetooth"
)

const serviceUUID = "4fafc201-1fb5-459e-8fcc-c5c9c331914b"
const dataCharacteristicUUID = "1c95d5e3-d8f7-413a-bf3d-7a2e5d7be87e"
const infoCharacteristicUUID = "beb5483e-36e1-4688-b7f5-ea07361b26a8"

func findDevice(a *bluetooth.Adapter) (bluetooth.ScanResult, error) {
	ch := make(chan bluetooth.ScanResult)
	var e error
	go func() {
		err := a.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.LocalName() == "SCRIBE" {
				adapter.StopScan()
				ch <- device
			}
		})
		if err != nil {
			e = err
		}
	}()

	select {
	case device := <-ch:
		return device, nil
	case <-time.After(1 * time.Second):
		e = errors.New("Timeout for finding device exceeded")
		return bluetooth.ScanResult{}, e
	}

}

func ConnectDevice(a *bluetooth.Adapter) (bluetooth.Device, bluetooth.DeviceCharacteristic, bluetooth.DeviceCharacteristic, error) {
	var device bluetooth.Device
	var dataChar, infoChar bluetooth.DeviceCharacteristic

	if a == nil {
		return device, dataChar, infoChar, errors.New("No ble adapter found")
	}

	res, err := findDevice(a)
	if err != nil {
		return device, dataChar, infoChar, err
	}

	device, err = a.Connect(res.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return device, dataChar, infoChar, err
	}

	serviceID, err := bluetooth.ParseUUID(serviceUUID)
	if err != nil {
		return device, dataChar, infoChar, err
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{serviceID})
	if err != nil {
		return device, dataChar, infoChar, err
	}

	if len(services) == 0 {
		return device, dataChar, infoChar, err
	}
	service := services[0]

	dataCharUUID, err := bluetooth.ParseUUID(dataCharacteristicUUID)
	if err != nil {
		return device, dataChar, infoChar, err
	}

	infoCharUUID, err := bluetooth.ParseUUID(infoCharacteristicUUID)
	if err != nil {
		return device, dataChar, infoChar, err
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{dataCharUUID, infoCharUUID})
	if err != nil {
		return device, dataChar, infoChar, err
	}

	if len(chars) != 2 {
		return device, dataChar, infoChar, err
	}

	if chars[0].UUID().String() == dataCharacteristicUUID {
		dataChar = chars[0]
		infoChar = chars[1]
	} else {
		dataChar = chars[1]
		infoChar = chars[0]
	}

	return device, dataChar, infoChar, nil
}

func DisconnectDevice(d bluetooth.Device) error {
	err := d.Disconnect()
	return err
}
