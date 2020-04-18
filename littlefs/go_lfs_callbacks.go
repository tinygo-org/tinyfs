package lfs

import (
	"fmt"
	"unsafe"

	"github.com/bgould/tinyfs"

	gopointer "github.com/mattn/go-pointer"
)

import "C"

const (
	debug bool = false
)

//export go_lfs_block_device_read
func go_lfs_block_device_read(ctx unsafe.Pointer, block uint32, offset uint32, buf unsafe.Pointer, size int) int {
	if debug {
		fmt.Printf("go_lfs_block_device_read: %v, %v, %v, %v, %v\n", ctx, block, offset, buf, size)
	}
	fs := restore(ctx)
	addr := fs.blockSize()*block + offset
	buffer := (*[1 << 28]byte)(buf)[:size:size]
	_, err := fs.dev.ReadAt(buffer, int64(addr))
	return go_lfs_block_errval("read", err)
}

//export go_lfs_block_device_prog
func go_lfs_block_device_prog(ctx unsafe.Pointer, block uint32, offset uint32, buf unsafe.Pointer, size int) int {
	if debug {
		fmt.Printf("go_lfs_block_device_prog: %v, %v, %v, %v, %v\n", ctx, block, offset, buf, size)
	}
	fs := restore(ctx)
	addr := fs.blockSize()*block + offset
	buffer := (*[1 << 28]byte)(buf)[:size:size]
	_, err := fs.dev.WriteAt(buffer, int64(addr))
	return go_lfs_block_errval("program", err)
}

//export go_lfs_block_device_erase
func go_lfs_block_device_erase(ctx unsafe.Pointer, block uint32) int {
	if debug {
		fmt.Printf("go_lfs_block_device_erase: %v, %v\n", ctx, block)
	}
	return go_lfs_block_errval("erase", restore(ctx).dev.EraseBlocks(int64(block), 1))
}

//export go_lfs_block_device_sync
func go_lfs_block_device_sync(ctx unsafe.Pointer) int {
	if debug {
		fmt.Printf("go_lfs_block_device_sync: %v\n", ctx)
	}
	fs := restore(ctx)
	if syncer, ok := fs.dev.(tinyfs.Syncer); ok {
		return go_lfs_block_errval("sync", syncer.Sync())
	}
	return errOK
}

func go_lfs_block_errval(op string, err error) int {
	if err != nil {
		if debug {
			println(op, "error:", err)
		}
		return int(errIO)
	}
	return errOK
}

func restore(ptr unsafe.Pointer) *LFS {
	return gopointer.Restore(ptr).(*LFS)
}
