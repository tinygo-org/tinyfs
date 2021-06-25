// +build tinygo

package main

import (
	"machine"
	"time"

	"tinygo.org/x/tinyfs/examples/console"
	"tinygo.org/x/tinyfs/littlefs"
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

	filesystem = littlefs.New(blockDevice)
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

	// Configure littlefs with parameters for caches and wear levelling
	filesystem.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})

	console.RunFor(blockDevice, filesystem)
}
