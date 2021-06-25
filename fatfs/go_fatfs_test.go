package fatfs

import (
	"testing"

	"tinygo.org/x/tinyfs"
)

const (
	testPageSize   = 64
	testBlockSize  = 256
	testBlockCount = 4096
)

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
