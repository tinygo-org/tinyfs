//go:build tinygo
// +build tinygo

package main

import (
	"machine"

	"tinygo.org/x/tinyfs/examples/console"
	"tinygo.org/x/tinyfs/littlefs"
)

var (
	blockDevice = machine.Flash
	filesystem  = littlefs.New(blockDevice)
)

func main() {
	// Configure littlefs with parameters for caches and wear levelling
	filesystem.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})
	console.RunFor(blockDevice, filesystem)
}
