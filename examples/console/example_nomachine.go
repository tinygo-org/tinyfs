//go:build !(nrf || nrf51 || nrf52 || nrf528xx || stm32f4 || stm32l4 || stm32wlx || atsamd21 || atsamd51 || atsame5x || rp2040)

package console

func lsblk_machine() {
	println("Unknown device")
}
