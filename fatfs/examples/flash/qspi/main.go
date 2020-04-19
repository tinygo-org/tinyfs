package main

import (
	"machine"

	example "github.com/bgould/go-fatfs/examples/flash"
	"tinygo.org/x/drivers/flash"
)

func main() {
	example.RunFor(
		flash.NewQSPI(
			machine.QSPI_CS,
			machine.QSPI_SCK,
			machine.QSPI_DATA0,
			machine.QSPI_DATA1,
			machine.QSPI_DATA2,
			machine.QSPI_DATA3,
		),
	)
}
