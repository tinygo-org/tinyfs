package littlefs

// #include <string.h>
// #include <stdlib.h>
// #include "./go_lfs.h"
import "C"

import (
	"errors"
	"io"
	"os"
	"time"
	"unsafe"

	"tinygo.org/x/tinyfs"
	"tinygo.org/x/tinyfs/internal/gopointer"
	"tinygo.org/x/tinyfs/internal/util"
)

const (
	errOK                 = C.LFS_ERR_OK          // No error
	errIO           Error = C.LFS_ERR_IO          // Error during device operation
	errCorrupt      Error = C.LFS_ERR_CORRUPT     // Corrupted
	errNoEntry      Error = C.LFS_ERR_NOENT       // No directory entry
	errEntryExists  Error = C.LFS_ERR_EXIST       // Entry already exists
	errNotDir       Error = C.LFS_ERR_NOTDIR      // Entry is not a dir
	errIsDir        Error = C.LFS_ERR_ISDIR       // Entry is a dir
	errDirNotEmpty  Error = C.LFS_ERR_NOTEMPTY    // Dir is not empty
	errBadFileNum   Error = C.LFS_ERR_BADF        // Bad file number
	errFileTooLarge Error = C.LFS_ERR_FBIG        // File too large
	errInvalidParam Error = C.LFS_ERR_INVAL       // Invalid parameter
	errNoSpace      Error = C.LFS_ERR_NOSPC       // No space left on device
	errNoMemory     Error = C.LFS_ERR_NOMEM       // No more memory available
	errNoAttr       Error = C.LFS_ERR_NOATTR      // No data/attr available
	errNameTooLong  Error = C.LFS_ERR_NAMETOOLONG // File name too long

	fileTypeReg fileType = C.LFS_TYPE_REG
	fileTypeDir fileType = C.LFS_TYPE_DIR
)

func translateFlags(osFlags int) C.int {
	var result C.int
	if osFlags&os.O_RDONLY > 0 {
		result |= C.LFS_O_RDONLY
	}
	if osFlags&os.O_WRONLY > 0 {
		result |= C.LFS_O_WRONLY
	}
	if osFlags&os.O_RDWR > 0 {
		result |= C.LFS_O_RDWR
	}
	if osFlags&os.O_CREATE > 0 {
		result |= C.LFS_O_CREAT
	}
	if osFlags&os.O_EXCL > 0 {
		result |= C.LFS_O_EXCL
	}
	if osFlags&os.O_TRUNC > 0 {
		result |= C.LFS_O_TRUNC
	}
	if osFlags&os.O_APPEND > 0 {
		result |= C.LFS_O_APPEND
	}
	return result
}

type fileType uint

type Error int

func (err Error) Error() string {
	switch err {
	case errIO:
		return "littlefs: Error during device operation"
	case errCorrupt:
		return "littlefs: Corrupted"
	case errNoEntry:
		return "littlefs: No directory entry"
	case errEntryExists:
		return "littlefs: Entry already exists"
	case errNotDir:
		return "littlefs: Entry is not a dir"
	case errIsDir:
		return "littlefs: Entry is a dir"
	case errDirNotEmpty:
		return "littlefs: Dir is not empty"
	case errBadFileNum:
		return "littlefs: Bad file number"
	case errFileTooLarge:
		return "littlefs: File too large"
	case errInvalidParam:
		return "littlefs: Invalid parameter"
	case errNoSpace:
		return "littlefs: No space left on device"
	case errNoMemory:
		return "littlefs: No more memory available"
	case errNoAttr:
		return "littlefs: No data/attr available"
	case errNameTooLong:
		return "littlefs: File name too long"
	default:
		return "littlefs: Unknown error"
	}
}

type Config struct {
	//ReadSize      uint32
	//ProgSize      uint32
	//BlockSize     uint32
	//BlockCount    uint32
	CacheSize     uint32
	LookaheadSize uint32
	BlockCycles   int32
}

type Info struct {
	ftyp fileType
	size uint32
	name string
}

func (info Info) Name() string {
	return info.name
}

func (info Info) Size() int64 {
	return int64(info.size)
}

func (info Info) IsDir() bool {
	return info.ftyp == fileTypeDir
}

func (info Info) Sys() interface{} {
	return nil
}

func (info Info) Mode() os.FileMode {
	v := os.FileMode(0777)
	if info.IsDir() {
		v |= os.ModeDir
	}
	return v
}

func (info Info) ModTime() time.Time {
	return time.Time{}
}

type LFS struct {
	dev tinyfs.BlockDevice
	ptr unsafe.Pointer
	lfs *C.struct_lfs
	cfg *C.struct_lfs_config
}

func New(blockdev tinyfs.BlockDevice) *LFS {
	return &LFS{
		dev: blockdev,
	}
}

func (l *LFS) Configure(config *Config) *LFS {
	l.lfs = C.go_lfs_new_lfs()
	l.cfg = C.go_lfs_new_lfs_config()
	*l.cfg = C.struct_lfs_config{
		context:        gopointer.Save(l),
		read_size:      C.lfs_size_t(l.dev.WriteBlockSize()),
		prog_size:      C.lfs_size_t(l.dev.WriteBlockSize()),
		block_size:     C.lfs_size_t(l.dev.EraseBlockSize()),
		block_count:    C.lfs_size_t(l.dev.Size() / l.dev.EraseBlockSize()),
		cache_size:     C.lfs_size_t(config.CacheSize),
		lookahead_size: C.lfs_size_t(config.LookaheadSize),
		block_cycles:   C.int32_t(config.BlockCycles),
	}
	C.go_lfs_set_callbacks(l.cfg)
	return l
}

func (l *LFS) blockSize() uint32 {
	return uint32(l.cfg.block_size)
}

func (l *LFS) Mount() error {
	return errval(C.lfs_mount(l.lfs, l.cfg))
}

func (l *LFS) Format() error {
	return errval(C.lfs_format(l.lfs, l.cfg))
}

func (l *LFS) Unmount() error {
	return errval(C.lfs_unmount(l.lfs))
}

func (l *LFS) Remove(path string) error {
	cs := cstring(path)
	defer C.free(unsafe.Pointer(cs))
	return errval(C.lfs_remove(l.lfs, cs))
}

func (l *LFS) Rename(oldPath string, newPath string) error {
	cs1, cs2 := cstring(oldPath), cstring(newPath)
	defer C.free(unsafe.Pointer(cs1))
	defer C.free(unsafe.Pointer(cs2))
	return errval(C.lfs_rename(l.lfs, cs1, cs2))
}

func (l *LFS) Stat(path string) (os.FileInfo, error) {
	cs := cstring(path)
	defer C.free(unsafe.Pointer(cs))
	info := C.struct_lfs_info{}
	if err := errval(C.lfs_stat(l.lfs, cs, &info)); err != nil {
		return nil, err
	}
	return &Info{
		ftyp: fileType(info._type),
		size: uint32(info.size),
		name: gostring(&info.name[0]),
	}, nil
}

func (l *LFS) Mkdir(path string, _ os.FileMode) error {
	cs := (*C.char)(cstring(path))
	defer C.free(unsafe.Pointer(cs))
	return errval(C.lfs_mkdir(l.lfs, cs))
}

func (l *LFS) Open(path string) (tinyfs.File, error) {
	return l.OpenFile(path, os.O_RDONLY)
}

func (l *LFS) OpenFile(path string, flags int) (tinyfs.File, error) {

	cs := (*C.char)(cstring(path))
	defer C.free(unsafe.Pointer(cs))
	file := &File{lfs: l, info: Info{name: path}}

	info := C.struct_lfs_info{}
	if err := errval(C.lfs_stat(l.lfs, cs, &info)); err == nil {
		file.info.ftyp = fileType(info._type)
		file.info.size = uint32(info.size)
	}

	var errno C.int
	if file.info.ftyp == fileTypeDir {
		file.typ = fileTypeDir
		file.hndl = unsafe.Pointer(C.go_lfs_new_lfs_dir())
		errno = C.lfs_dir_open(l.lfs, file.dirptr(), cs)
	} else {
		file.typ = fileTypeReg
		file.hndl = unsafe.Pointer(C.go_lfs_new_lfs_file())
		errno = C.lfs_file_open(l.lfs, file.fileptr(), cs, C.int(translateFlags(flags)))
	}

	if err := errval(errno); err != nil {
		if file.hndl != nil {
			C.free(file.hndl)
			file.hndl = nil
		}
		return nil, err
	}

	return file, nil
}

// Size finds the current size of the filesystem
//
// Note: Result is best effort. If files share COW structures, the returned
// size may be larger than the filesystem actually is.
//
// Returns the number of allocated blocks, or a negative error code on failure.
func (l *LFS) Size() (n int, err error) {
	errno := C.int(C.lfs_fs_size(l.lfs))
	if errno < 0 {
		return 0, errval(errno)
	}
	return int(errno), nil
}

type File struct {
	lfs  *LFS
	typ  fileType
	hndl unsafe.Pointer
	info Info
}

func (f *File) dirptr() *C.struct_lfs_dir {
	return (*C.struct_lfs_dir)(f.hndl)
}

func (f *File) fileptr() *C.struct_lfs_file {
	return (*C.struct_lfs_file)(f.hndl)
}

// Name returns the name of the file as presented to OpenFile
func (f *File) Name() string {
	return f.info.name
}

// Close the file; any pending writes are written out to storage
func (f *File) Close() error {
	if f.hndl != nil {
		defer func() {
			C.free(f.hndl)
			f.hndl = nil
		}()
		switch f.typ {
		case fileTypeReg:
			return errval(C.lfs_file_close(f.lfs.lfs, f.fileptr()))
		case fileTypeDir:
			return errval(C.lfs_dir_close(f.lfs.lfs, f.dirptr()))
		default:
			panic("lfs: unknown typ for file handle")
		}
	}
	return nil
}

func (f *File) Read(buf []byte) (n int, err error) {
	if f.IsDir() {
		return 0, errIsDir
	}
	bufptr := unsafe.Pointer(&buf[0])
	buflen := C.lfs_size_t(len(buf))
	errno := C.int(C.lfs_file_read(f.lfs.lfs, f.fileptr(), bufptr, buflen))
	if errno > 0 {
		return int(errno), nil
	} else if errno == 0 {
		// TODO: any extra checks needed here?
		return 0, io.EOF
	} else {
		return 0, errval(errno)
	}
}

// Seek changes the position of the file
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	errno := C.int(C.lfs_file_seek(f.lfs.lfs, f.fileptr(), C.lfs_soff_t(offset), C.int(whence)))
	if errno < 0 {
		return -1, errval(errno)
	}
	return int64(errno), nil
}

// Tell returns the position of the file
func (f *File) Tell() (ret int64, err error) {
	errno := C.int(C.lfs_file_tell(f.lfs.lfs, f.fileptr()))
	if errno < 0 {
		return -1, errval(errno)
	}
	return int64(errno), nil
}

// Rewind changes the position of the file to the beginning of the file
func (f *File) Rewind() (err error) {
	return errval(C.lfs_file_rewind(f.lfs.lfs, f.fileptr()))
}

// Size returns the size of the file
func (f *File) Size() (int64, error) {
	errno := C.int(C.lfs_file_size(f.lfs.lfs, f.fileptr()))
	if errno < 0 {
		return -1, errval(errno)
	}
	return int64(errno), nil
}

// Stat satisfies the `fs.File` interface
func (f *File) Stat() (os.FileInfo, error) {
	return f.info, nil
}

// Sync synchronizes to storage so that any pending writes are written out.
func (f *File) Sync() error {
	return errval(C.lfs_file_sync(f.lfs.lfs, f.fileptr()))
}

// Truncate the size of the file to the specified size
func (f *File) Truncate(size uint32) error {
	return errval(C.lfs_file_truncate(f.lfs.lfs, f.fileptr(), C.lfs_off_t(size)))
}

func (f *File) Write(buf []byte) (n int, err error) {
	bufptr := unsafe.Pointer(&buf[0])
	buflen := C.lfs_size_t(len(buf))
	errno := C.lfs_file_write(f.lfs.lfs, f.fileptr(), bufptr, buflen)
	if errno > 0 {
		return int(errno), nil
	} else {
		return 0, errval(C.int(errno))
	}
}

func (f *File) IsDir() bool {
	return f.typ == fileTypeDir
}

func (f *File) Readdir(n int) (infos []os.FileInfo, err error) {
	if n > 0 {
		return nil, errors.New("n > 0 is not supported yet")
	}
	if !f.IsDir() {
		return nil, errNotDir
	}
	for {
		var info C.struct_lfs_info
		i := C.lfs_dir_read(f.lfs.lfs, f.dirptr(), &info)
		if i == 0 {
			return
		}
		if i < 0 {
			err = errval(C.int(i))
			return
		}
		name := gostring(&info.name[0])
		if name == "." || name == ".." {
			continue // littlefs returns . and .., but Readdir() in Go does not
		}
		infos = append(infos, Info{
			ftyp: fileType(info._type),
			size: uint32(info.size),
			name: name,
		})
	}
}

func errval(errno C.int) error {
	if errno < errOK {
		return Error(errno)
	}
	return nil
}

func cstring(s string) *C.char {
	return (*C.char)(util.CString(s))
}

func gostring(s *C.char) string {
	return util.GoString(unsafe.Pointer(s))
}
