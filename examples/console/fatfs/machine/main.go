//go:build tinygo
// +build tinygo

package main

import (
	"machine"

	"tinygo.org/x/tinyfs/examples/console"
	"tinygo.org/x/tinyfs/fatfs"
)

var (
	blockDevice = machine.Flash
	filesystem  = fatfs.New(blockDevice)
)

func main() {
	// Configure FATFS with sector size (must match value in ff.h - use 512)
	filesystem.Configure(&fatfs.Config{
		SectorSize: 512,
	})

	console.RunFor(blockDevice, filesystem)
}
