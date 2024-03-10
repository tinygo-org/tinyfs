package tinyfs

import (
	"fmt"
	"io"
	"os"
)

type Filesystem interface {
	Format() error
	Mkdir(path string, _ os.FileMode) error
	Mount() error
	Open(path string) (File, error)
	OpenFile(path string, flags int) (File, error)
	Remove(path string) error
	Rename(oldPath string, newPath string) error
	Stat(path string) (os.FileInfo, error)
	Unmount() error
}

// File specifies the common behavior of the file abstraction in TinyFS; this
// interface may be changed or superseded by the TinyGo os.FileHandle interface
// if/when that is merged and standardized.
type File interface {
	FileHandle
	io.Seeker
	IsDir() bool
	Readdir(n int) (infos []os.FileInfo, err error)
	Stat() (info os.FileInfo, err error)
}

// FileHandle is a copy of the experimental os.FileHandle interface in TinyGo
type FileHandle interface {
	// Read reads up to len(b) bytes from the file.
	Read(b []byte) (n int, err error)

	// Write writes up to len(b) bytes to the file.
	Write(b []byte) (n int, err error)

	// Close closes the file, making it unusable for further writes.
	Close() (err error)
}

// A BlockDevice is the raw device that is meant to store a filesystem.
type BlockDevice interface {

	// ReadAt reads the given number of bytes from the block device.
	io.ReaderAt

	// WriteAt writes the given number of bytes to the block device.
	io.WriterAt

	// Size returns the number of bytes in this block device.
	Size() int64

	// WriteBlockSize returns the block size in which data can be written to
	// memory. It can be used by a client to optimize writes, non-aligned writes
	// should always work correctly.
	WriteBlockSize() int64

	// EraseBlockSize returns the smallest erasable area on this particular chip
	// in bytes. This is used for the block size in EraseBlocks.
	// It must be a power of two, and may be as small as 1. A typical size is 4096.
	EraseBlockSize() int64

	// EraseBlocks erases the given number of blocks. An implementation may
	// transparently coalesce ranges of blocks into larger bundles if the chip
	// supports this. The start and len parameters are in block numbers, use
	// EraseBlockSize to map addresses to blocks.
	EraseBlocks(start, len int64) error
}

type Syncer interface {
	// Sync triggers the devices to commit any pending or cached operations
	Sync() error
}

// MemBlockDevice is a block device implementation backed by a byte slice
type MemBlockDevice struct {
	memory     []byte
	blankBlock []byte
	blockCount uint32
	blockSize  uint32
	pageSize   uint32
}

var _ BlockDevice = (*MemBlockDevice)(nil)

func NewMemoryDevice(pageSize, blockSize, blockCount int) *MemBlockDevice {
	dev := &MemBlockDevice{
		memory:     make([]byte, blockSize*blockCount),
		blankBlock: make([]byte, blockSize),
		pageSize:   uint32(pageSize),
		blockSize:  uint32(blockSize),
		blockCount: uint32(blockCount),
	}
	for i := range dev.blankBlock {
		dev.blankBlock[i] = 0xff
	}
	for i := uint32(0); i < uint32(blockCount); i++ {
		if err := dev.eraseBlock(i); err != nil {
			panic(fmt.Sprintf("could not initialize block %d: %s", i, err.Error()))
		}
	}
	return dev
}

func (bd *MemBlockDevice) ReadAt(buf []byte, off int64) (n int, err error) {
	return copy(buf, bd.memory[off:]), nil
}

func (bd *MemBlockDevice) WriteAt(buf []byte, off int64) (n int, err error) {
	return copy(bd.memory[off:], buf), nil
}

func (bd *MemBlockDevice) Size() int64 {
	return int64(bd.blockSize * bd.blockCount)
}

func (bd *MemBlockDevice) WriteBlockSize() int64 {
	return int64(bd.pageSize)
}

func (bd *MemBlockDevice) EraseBlockSize() int64 {
	return int64(bd.blockSize)
}

func (bd *MemBlockDevice) EraseBlocks(start int64, len int64) error {
	for i := int64(0); i < len; i++ {
		if err := bd.eraseBlock(uint32(start + i)); err != nil {
			return err
		}
	}
	return nil
}

func (bd *MemBlockDevice) eraseBlock(block uint32) error {
	copy(bd.memory[bd.blockSize*block:], bd.blankBlock)
	return nil
}

func (bd *MemBlockDevice) Sync() error {
	return nil
}
