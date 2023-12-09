//go:build !tinygo
// +build !tinygo

package main

import (
	"os"

	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/fatfs"
)

// FileBlockDevice is a block device implementation backed by a byte slice
type FileBlockDevice struct {
	file       *os.File
	blankBlock []byte
	pageSize   uint32
	blockSize  uint32
	blockCount uint32
}

var _ tinyfs.BlockDevice = (*FileBlockDevice)(nil)

func NewFileDevice(file *os.File, pageSize int, blockSize int, blockCount int) *FileBlockDevice {
	dev := &FileBlockDevice{
		blankBlock: make([]byte, blockSize),
		pageSize:   uint32(pageSize),
		blockSize:  uint32(blockSize),
		blockCount: uint32(blockCount),
	}
	return dev
}

func (bd *FileBlockDevice) ReadAt(buf []byte, off int64) (n int, err error) {
	return bd.file.ReadAt(buf, off)
	//return copy(buf, bd.memory[off:]), nil
}

func (bd *FileBlockDevice) WriteAt(buf []byte, off int64) (n int, err error) {
	return bd.file.WriteAt(buf, off)
}

func (bd *FileBlockDevice) Size() int64 {
	return int64(bd.blockSize * bd.blockCount)
}

func (bd *FileBlockDevice) SectorSize() int {
	return fatfs.SectorSize
}

func (bd *FileBlockDevice) WriteBlockSize() int64 {
	return int64(bd.pageSize)
}

func (bd *FileBlockDevice) EraseBlockSize() int64 {
	return int64(bd.blockSize)
}

func (bd *FileBlockDevice) EraseBlocks(start int64, len int64) error {
	for i := int64(0); i < len; i++ {
		if err := bd.eraseBlock(uint32(start + i)); err != nil {
			return err
		}
	}
	return nil
}

func (bd *FileBlockDevice) eraseBlock(block uint32) error {
	_, err := bd.file.WriteAt(bd.blankBlock, int64(bd.blockSize*block))
	return err
}

func (bd *FileBlockDevice) Sync() error {
	return nil
}
