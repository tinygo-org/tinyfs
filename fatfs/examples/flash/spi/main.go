package main

import (
	"machine"

	example "github.com/bgould/go-fatfs/examples/flash"

	"tinygo.org/x/drivers/flash"
)

func main() {
	example.RunFor(
		flash.NewSPI(
			&machine.SPI1,
			machine.SPI1_MOSI_PIN,
			machine.SPI1_MISO_PIN,
			machine.SPI1_SCK_PIN,
			machine.SPI1_CS_PIN,
		),
	)
}
