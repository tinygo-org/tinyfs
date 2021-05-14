// +build tinygo

package main

import (
	"machine"
	"time"

	"github.com/tinygo-org/tinyfs/examples/console"
	"github.com/tinygo-org/tinyfs/fatfs"
	"tinygo.org/x/drivers/flash"
)

var (
	blockDevice = flash.NewSPI(
		&machine.SPI1,
		machine.SPI1_MOSI_PIN,
		machine.SPI1_MISO_PIN,
		machine.SPI1_SCK_PIN,
		machine.SPI1_CS_PIN,
	)

	filesystem = fatfs.New(blockDevice)
)

func main() {

	// Configure the flash device using the default auto-identifier function
	config := &flash.DeviceConfig{Identifier: flash.DefaultDeviceIdentifier}
	if err := blockDevice.Configure(config); err != nil {
		for {
			time.Sleep(5 * time.Second)
			println("Config was not valid: "+err.Error(), "\r")
		}
	}

	// Configure FATFS with sector size (must match value in ff.h - use 512)
	filesystem.Configure(&fatfs.Config{
		SectorSize: 512,
	})

	console.RunFor(blockDevice, filesystem)
}
