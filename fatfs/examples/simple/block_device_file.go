// +build !tinygo

package main

import (
	"os"
)

// FileBlockDevice is a block device implementation backed by a byte slice
type FileBlockDevice struct {
	file       *os.File
	blankBlock []byte
	blockSize  uint32
	blockCount uint32
}

var _ BlockDevice = (*FileBlockDevice)(nil)

func NewFileDevice(file *os.File, blockSize int, blockCount int) *FileBlockDevice {
	dev := &FileBlockDevice{
		blankBlock: make([]byte, blockSize),
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
	return SectorSize
}

func (bd *FileBlockDevice) EraseBlockSize() int {
	return int(bd.blockSize)
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
