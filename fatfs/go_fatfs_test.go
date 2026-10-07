//go:build !tinygo

package fatfs

import (
	"encoding/binary"
	"os"
	"testing"
	"time"

	"tinygo.org/x/tinyfs"
)

const (
	testPageSize   = 64
	testBlockSize  = 256
	testBlockCount = 4096
)

func TestStoredModificationTime(t *testing.T) {
	fs, dev, unmount := createTestFS(t)
	defer unmount()
	if typ := Type(fs.fs.fs_type); typ != TypeFAT12 && typ != TypeFAT16 {
		t.Fatalf("expected FAT12 or FAT16, got %s", typ)
	}
	before := time.Now().Truncate(2 * time.Second)
	f, err := fs.OpenFile("TIME.TXT", os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	check(t, err)
	_, err = f.Write([]byte("timestamp"))
	check(t, err)
	check(t, f.Close())
	after := time.Now()

	root := make([]byte, int(fs.fs.n_rootdir)*32)
	_, err = dev.ReadAt(root, int64(fs.fs.dirbase)*SectorSize)
	check(t, err)
	for offset := 0; offset < len(root); offset += 32 {
		entry := root[offset : offset+32]
		if string(entry[:11]) != "TIME    TXT" {
			continue
		}
		// FAT directory entries store modification time at bytes 22-25.
		clock := binary.LittleEndian.Uint16(entry[22:24])
		date := binary.LittleEndian.Uint16(entry[24:26])
		got := time.Date(1980+int(date>>9), time.Month(date>>5&15), int(date&31),
			int(clock>>11), int(clock>>5&63), int(clock&31)*2, 0, time.Local)
		if got.Before(before) || got.After(after) {
			t.Fatalf("stored modification time = %s, want between %s and %s", got, before, after)
		}
		return
	}
	t.Fatal("TIME.TXT not found in root directory")
}

func TestType_String(t *testing.T) {
	expectString(t, "fatfs: (1) A hard error occurred in the low level disk I/O layer", FileResultErr.Error())
}

func TestFormat(t *testing.T) {
	//dev2 := NewMemoryDevice(4096, 64)
	//fs2 := New(dev2)
	t.Run("BasicMounting", func(t *testing.T) {
		// file, err := os.Open("trinket_filesystem.img")
		// check(t, err)
		// dev := NewFileDevice(file, 512, 29)
		fs, _, umount := createTestFS(t)
		defer umount()
		n, err := fs.Free()
		check(t, err)
		println("free:", n)
	})
}

func createTestFS(t *testing.T) (*FATFS, tinyfs.BlockDevice, func()) {
	// create/format/mount the filesystem
	dev := tinyfs.NewMemoryDevice(testPageSize, testBlockSize, testBlockCount)
	fs := New(dev)
	println("formatting")
	fs.Configure(&Config{SectorSize: 512})
	//check(t, fs.Configure())
	if err := fs.Format(); err != nil {
		t.Fatal(err)
	}
	/*
		buf := make([]byte, 512)
		dev.ReadAt(buf, 0)
		xxdfprint(os.Stdout, 0, buf)
		println("format successful; mounting")
	*/
	if err := fs.Mount(); err != nil {
		t.Error("Could not mount", err)
	}
	return fs, dev, func() {
		//if err := fs.Unmount(); err != nil {
		//	t.Error("Could not ummount", err)
		//}
	}
}

func TestDirectories(t *testing.T) {

	const (
		largeSize = 128
	)
	t.Run("RootDirectory", func(t *testing.T) {
		fs, _, unmount := createTestFS(t)
		defer unmount()
		f, err := fs.Open("/")
		check(t, err)
		check(t, f.Close())
	})
	t.Run("DirectoryCreation", func(t *testing.T) {
		fs, _, unmount := createTestFS(t)
		defer unmount()
		check(t, fs.Mkdir("potato", 0777))
		info, err := fs.Stat("potato")
		check(t, err)
		println("potato: ", info.Name(), info.IsDir(), info.Mode())
	})
}

func TestFiles(t *testing.T) {
	t.Run("BasicFile", func(t *testing.T) {
		fs, _, unmount := createTestFS(t)
		defer unmount()
		f, err := fs.OpenFile("hello_world.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
		check(t, err)
		n, err := f.Write([]byte("hello world"))
		check(t, err)
		check(t, f.Close())
		if n != 11 {
			t.Fail()
		}
		f, err = fs.Open("hello_world.txt")
		check(t, err)
		info, err := f.Stat()
		check(t, err)
		if info.IsDir() {
			t.Fail()
		}
		if info.Name() != "hello_world.txt" {
			t.Fail()
		}
		if info.Size() != 11 {
			t.Fail()
		}
	})
}

func expectString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatalf("expected \"%s\", was actually \"%s\"", expected, actual)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
