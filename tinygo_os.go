//go:build tinygo && osfilesystem
// +build tinygo,osfilesystem

package tinyfs

import "os"

type OSFilesystem struct {
	FS Filesystem
}

// OpenFile opens the named file.
func (fs *OSFilesystem) OpenFile(name string, flag int, perm os.FileMode) (os.FileHandle, error) {
	return fs.FS.OpenFile(name, flag)
}

// Mkdir creates a new directory with the specified permission (before
// umask). Some filesystems may not support directories or permissions.
func (fs *OSFilesystem) Mkdir(name string, perm os.FileMode) error {
	return fs.FS.Mkdir(name, perm)
}

// Remove removes the named file or (empty) directory.
func (fs *OSFilesystem) Remove(name string) error {
	return fs.FS.Remove(name)
}

var _ os.Filesystem = (*OSFilesystem)(nil)
