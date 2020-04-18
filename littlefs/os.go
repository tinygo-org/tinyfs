// +build tinygo

package lfs

import "os"

type Filesystem struct {
	LFS *LFS
}

// OpenFile opens the named file.
func (fs *Filesystem) OpenFile(name string, flag int, perm os.FileMode) (os.FileHandle, error) {
	return fs.LFS.OpenFile(name, flag)
}

// Mkdir creates a new directory with the specified permission (before
// umask). Some filesystems may not support directories or permissions.
func (fs *Filesystem) Mkdir(name string, perm os.FileMode) error {
	return fs.LFS.Mkdir(name, perm)
}

// Remove removes the named file or (empty) directory.
func (fs *Filesystem) Remove(name string) error {
	return fs.LFS.Remove(name)
}

var _ os.Filesystem = (*Filesystem)(nil)
