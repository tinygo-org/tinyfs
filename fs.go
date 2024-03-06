package tinyfs

import (
	"io/fs"
)

func NewTinyFS(filesystem Filesystem) fs.FS {
	return &fsWrapper{fs: filesystem}
}

type fsWrapper struct {
	fs Filesystem
}

func (w *fsWrapper) Open(name string) (f fs.File, err error) {
	return w.fs.Open(name)
}

// type assertion to ensure tinyfs.File satisfies fs.File
var _ fs.File = (File)(nil)
