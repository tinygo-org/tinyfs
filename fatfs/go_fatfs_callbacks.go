package fatfs

// #include "diskio.h"
// #include "ff.h"
import "C"

import (
	"time"
	"unsafe"

	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/internal/gopointer"
)

const (
	debug = false
)

//export go_fatfs_disk_read
func go_fatfs_disk_read(drv unsafe.Pointer, bufptr unsafe.Pointer, sector uint32, count uint) int {
	if debug {
		println("disk_read:", sector, count)
	}
	bdev := restore(drv).dev
	sect := int(SectorSize)
	size := sect * int(count)
	addr := int64(sector) * int64(sect)
	buffer := (*[1 << 28]byte)(bufptr)[:size:size]
	if _, err := bdev.ReadAt(buffer, addr); err != nil {
		//println("disk_read error:", err)
		return C.RES_ERROR
	}
	return C.RES_OK
}

//export go_fatfs_disk_write
func go_fatfs_disk_write(drv unsafe.Pointer, bufptr unsafe.Pointer, sector uint32, count uint) int {
	//println("disk_write:", sector, count)
	bdev := restore(drv).dev
	sect := int(SectorSize)
	size := sect * int(count)
	addr := int64(sector) * int64(sect)
	buffer := (*[1 << 28]byte)(bufptr)[:size:size]
	//xxdfprint(os.Stdout, 0, buffer)
	if _, err := bdev.WriteAt(buffer, addr); err != nil {
		//println("disk_write error:", err)
		return C.RES_ERROR
	}
	return C.RES_OK
}

//export go_fatfs_disk_ioctl
func go_fatfs_disk_ioctl(drv unsafe.Pointer, cmd uint8, param unsafe.Pointer) int {
	//println("disk_ioctl:", cmd)
	bdev := restore(drv).dev
	switch cmd {
	case C.CTRL_SYNC:
		// Complete pending write process (needed at _FS_READONLY == 0)
		if syncer, ok := bdev.(tinyfs.Syncer); ok {
			if err := syncer.Sync(); err != nil {
				return C.RES_ERROR
			}
		}
	case C.GET_SECTOR_COUNT:
		// Get media size (needed at _USE_MKFS == 1)
		*((*C.DWORD)(param)) = C.DWORD(bdev.Size() / SectorSize)
	case C.GET_SECTOR_SIZE:
		// Get sector size (needed at _MAX_SS != _MIN_SS)
		*((*C.WORD)(param)) = C.WORD(SectorSize)
	case C.GET_BLOCK_SIZE:
		// Get erase block size (needed at _USE_MKFS == 1)
		// FIXME: not really sure why this doesn't work
		*((*C.DWORD)(param)) = C.DWORD(bdev.EraseBlockSize() / SectorSize)
	case C.IOCTL_INIT:
		// FIXME: not really sure what this would be used for
		*((*C.DSTATUS)(param)) = C.DSTATUS(0)
	case C.IOCTL_STATUS:
		// FIXME: not really sure what this would be used for
		*((*C.DSTATUS)(param)) = C.DSTATUS(0)
	}
	return C.RES_OK
}

//export go_fatfs_get_fattime
func go_fatfs_get_fattime() (t uint32) {
	const (
		epochYear      = 1980 // years are stored as an offset from 1980
		yearShift      = 25   // year offset, bits 25-31
		monthShift     = 21   // month, bits 21-24
		dayShift       = 16   // day of month, bits 16-20
		hourShift      = 11   // hour, bits 11-15
		minuteShift    = 5    // minute, bits 5-10
		monthMask      = 0x0f // month, 4 bits
		dayMask        = 0x1f // day of month, 5 bits
		hourMask       = 0x1f // hour, 5 bits
		minuteMask     = 0x3f // minute, 6 bits
		secondMask     = 0x1f // two-second units, 5 bits
		secondsPerUnit = 2    // FAT timestamps have two-second resolution
	)

	now := time.Now().UTC()
	year, month, day := now.Date()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	t |= uint32(year-epochYear) << yearShift
	t |= (uint32(month) & monthMask) << monthShift
	t |= (uint32(day) & dayMask) << dayShift
	t |= (uint32(hour) & hourMask) << hourShift
	t |= (uint32(minute) & minuteMask) << minuteShift
	t |= (uint32(second) / secondsPerUnit) & secondMask
	return t
}

func restore(ptr unsafe.Pointer) *FATFS {
	return gopointer.Restore(ptr).(*FATFS)
}
