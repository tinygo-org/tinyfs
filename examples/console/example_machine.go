//go:build nrf || nrf51 || nrf52 || nrf528xx || stm32f4 || stm32l4 || stm32wlx || atsamd21 || atsamd51 || atsame5x || rp2040

package console

import (
	"fmt"
	"machine"
)

func lsblk_machine() {
	if blockdev == machine.Flash {
		fmt.Printf(
			"\n-------------------------------------\r\n"+
				" Device Information:  \r\n"+
				"-------------------------------------\r\n"+
				" flash data start: %08X\r\n"+
				" flash data end:  %08X\r\n"+
				"-------------------------------------\r\n\r\n",
			machine.FlashDataStart(),
			machine.FlashDataEnd())
		return
	}
	println("Unknown device")
}
