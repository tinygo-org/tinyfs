package tinyfs

import (
	"errors"
	"math/bits"
)

var (
	ErrInvalidBlockSize = errors.New("tinyfs: write/erase block size too big or not aligned")
	ErrDeviceTooSmall   = errors.New("tinyfs: requested device size is smaller than underlying device")
)

type blockDevOffset struct {
	dev    BlockDevice
	offset int64
}

// Create a new virtual BlockDevice that is the upper part of the given full
// BlockDevice, with the given size. This can be useful for example to use only
// a fixed size machine.Flash object that won't change on firmware changes.
func SliceBlockDev(full BlockDevice, size int64) (BlockDevice, error) {
	if bits.OnesCount16(uint16(full.WriteBlockSize())) != 1 || bits.OnesCount16(uint16(full.EraseBlockSize())) != 1 {
		// Check that the write/erase block size is a power of two.
		// This is normally the case, and makes the following calculations
		// easier.
		return nil, ErrInvalidBlockSize
	}

	// Check whether the requested size can be used (is a multiple of the block
	// size).
	blockSize := max(full.EraseBlockSize(), full.WriteBlockSize())
	if size/blockSize*blockSize != size {
		return nil, ErrInvalidBlockSize
	}

	fullSize := full.Size()
	if fullSize < size {
		// Not big enough to provide the expected device size.
		return nil, ErrDeviceTooSmall
	}
	if fullSize == size {
		// Shortcut: the requested size is exactly big enough.
		return full, nil
	}
	offset := fullSize - size
	return &blockDevOffset{
		dev:    full,
		offset: offset,
	}, nil
}

func (b *blockDevOffset) ReadAt(buf []byte, offset int64) (n int, err error) {
	return b.dev.ReadAt(buf, offset+b.offset)
}

func (b *blockDevOffset) WriteAt(buf []byte, offset int64) (n int, err error) {
	return b.dev.WriteAt(buf, offset+b.offset)
}

func (b *blockDevOffset) Size() int64 {
	return b.dev.Size() - b.offset
}

func (b *blockDevOffset) EraseBlockSize() int64 {
	return b.dev.EraseBlockSize()
}

func (b *blockDevOffset) WriteBlockSize() int64 {
	return b.dev.WriteBlockSize()
}

func (b *blockDevOffset) EraseBlocks(start, len int64) error {
	blockOffset := b.offset / b.EraseBlockSize()
	return b.dev.EraseBlocks(start+blockOffset, len)
}
