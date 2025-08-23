package communication

import (
	"encoding/binary"
	"time"

	"tinygo.org/x/bluetooth"
)

const dataTransferVersion uint8 = 1

func SendData(dataChar, infoChar bluetooth.DeviceCharacteristic, data string) error {
	_, err := infoChar.WriteWithoutResponse(infoToBytes(dataTransferVersion, uint16(len(data)), 1))
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	chunkSize := 20

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize 
		if end > len(data) {
			end = len(data)
		}
		_, err = dataChar.WriteWithoutResponse([]byte(data[i:end]))
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)
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

func infoToBytes(v uint8, size uint16, font uint8) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(v)
	buf[1] = byte(font)
	binary.LittleEndian.PutUint16(buf[2:], size)

	return buf
}
