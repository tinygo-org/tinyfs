// +build tinygo

package main

import (
	"machine"

	example "github.com/bgould/tinyfs/examples/flash"
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
