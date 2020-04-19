package main

import (
	"machine"
	"time"

	"github.com/bgould/tinyfs/examples/console"
	"github.com/bgould/tinyfs/littlefs"
	"tinygo.org/x/drivers/flash"
)

var (
	blockDevice = flash.NewQSPI(
		machine.QSPI_CS,
		machine.QSPI_SCK,
		machine.QSPI_DATA0,
		machine.QSPI_DATA1,
		machine.QSPI_DATA2,
		machine.QSPI_DATA3,
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
