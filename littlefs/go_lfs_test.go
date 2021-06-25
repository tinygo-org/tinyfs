package littlefs

import (
	"io"
	"math/rand"
	"os"
	"testing"

	"tinygo.org/x/tinyfs"
)

const (
	testPageSize   = 64
	testBlockSize  = 256
	testBlockCount = 2048
)

var defaultConfig = &Config{
	//	ReadSize:      16,
	//	ProgSize:      16,
	//	BlockSize:     512,
	//	BlockCount:    1024,
	CacheSize:     128,
	LookaheadSize: 128,
	BlockCycles:   500,
}

var zeroBlock = make([]byte, testBlockSize)

func TestFormat(t *testing.T) {
	dev := tinyfs.NewMemoryDevice(testPageSize, testBlockSize, testBlockCount)
	fs := New(dev)
	fs.Configure(defaultConfig)

	t.Run("BasicFormatting", func(t *testing.T) {
		if err := fs.Format(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("BasicMounting", func(t *testing.T) {
		if err := fs.Mount(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("InvalidSuperblocks", func(t *testing.T) {
		if err := fs.Unmount(); err != nil {
			t.Fatal(err)
		}
		if _, err := dev.WriteAt(zeroBlock, 0); err != nil {
			t.Fatal(err)
		}
		if _, err := dev.WriteAt(zeroBlock, testBlockSize); err != nil {
			t.Fatal(err)
		}
		if err := fs.Format(); err == nil {
			t.Log("expected error when formatting")
			if err := fs.Mount(); err == nil {
				t.Log("expected error when mounting")
				//t.Fail()
			}
		} else if err != errNoSpace {
			t.Logf("expected error to be ErrNoSpace; was %s", err)
			t.Fail()
		}
	})

	t.Run("InvalidMount", func(t *testing.T) {
		dev := tinyfs.NewMemoryDevice(testPageSize, testBlockSize, testBlockCount)
		fs := New(dev)
		fs.Configure(defaultConfig)
		if err := fs.Mount(); err == nil {
			t.Log("expected error when mounting")
			t.Fail()
		} else if err != errCorrupt {
			t.Logf("expected error to be ErrCorrupt; was %s", err)
			t.Fail()
		}
	})

	t.Run("ExpandingSuperblock", func(t *testing.T) {
		if err := fs.Format(); err != nil {
			t.Fatal(err)
		}
		if err := fs.Mount(); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 100; i++ {
			if err := fs.Mkdir("dummy", 0777); err != nil {
				t.Fatalf("failing mkdir on iteration %d: %s", i, err)
			}
			if err := fs.Remove("dummy"); err != nil {
				t.Fatalf("failing remove on iteration %d: %s", i, err)
			}
		}
		if err := fs.Unmount(); err != nil {
			t.Fatal(err)
		}
		if err := fs.Mount(); err != nil {
			t.Fatalf("failed to re-mount: %s", err)
		}
		if err := fs.Mkdir("dummy", 0777); err != nil {
			t.Fatalf("fail to mkdir: %s", err)
		}
	})
}

func TestFiles(t *testing.T) {
	fs, _, unmount := createTestFS(t, defaultConfig)
	defer unmount()

	const (
		smallSize  = 32
		mediumSize = 8192
		largeSize  = 262144
	)

	t.Run("SimpleFileTest", func(t *testing.T) {
		f, err := fs.OpenFile("hello", os.O_WRONLY|os.O_CREATE)
		if err != nil {
			t.Error(err)
		}
		s := "Hello World!"
		if _, err := f.Write([]byte(s)); err != nil {
			t.Error(err)
		}
		if err := f.Close(); err != nil {
			t.Error(err)
		}
		f, err = fs.Open("hello")
		if err != nil {
			t.Error(err)
		}
		buf := make([]byte, 24)
		n, err := f.Read(buf)
		if err != nil {
			t.Error(err)
		}
		s2 := string(buf[:n])
		if n != len(s) {
			t.Fail()
		}
		if s != s2 {
			t.Fail()
		}
		if err := f.Close(); err != nil {
			t.Error(err)
		}
	})

	t.Run("SmallFile", func(t *testing.T) {
		writeFileTest(t, fs, smallSize, "smallavocado")
		readFileTest(t, fs, smallSize, "smallavocado")
	})

	t.Run("MediumFile", func(t *testing.T) {
		writeFileTest(t, fs, mediumSize, "mediumavocado")
		readFileTest(t, fs, mediumSize, "mediumavocado")
	})

	t.Run("LargeFile", func(t *testing.T) {
		writeFileTest(t, fs, largeSize, "largeavocado")
		readFileTest(t, fs, largeSize, "largeavocado")
	})

	t.Run("ZeroFile", func(t *testing.T) {
		writeFileTest(t, fs, 0, "noavocado")
		readFileTest(t, fs, 0, "noavocado")
	})
}

func TestDirectories(t *testing.T) {

	fs, _, unmount := createTestFS(t, defaultConfig)
	defer unmount()
	const (
		largeSize = 128
	)

	t.Run("RootDirectory", func(t *testing.T) {
		f, err := fs.Open("/")
		check(t, err)
		check(t, f.Close())
	})

	t.Run("DirectoryCreation", func(t *testing.T) {
		check(t, fs.Mkdir("potato", 0777))
	})

	t.Run("FileCreation", func(t *testing.T) {
		f, err := fs.OpenFile("burrito", os.O_CREATE|os.O_WRONLY)
		check(t, err)
		check(t, f.Close())
	})
	/*
		t.Run("DirectoryIteration", func(t *testing.T) {
			dir, err := fs.Open("/")
			check(t, err)
			defer check(t, dir.Close())
			infos, err := dir.Readdir(0)
			if err != nil {

			}
		})
	*/
	t.Run("DirectoryFailures", func(t *testing.T) {

	})

	t.Run("NestedDirectories", func(t *testing.T) {

	})

	t.Run("MultiBlockDirectory", func(t *testing.T) {

	})

	t.Run("DirectoryRemove", func(t *testing.T) {

	})

	t.Run("DirectoryRename", func(t *testing.T) {

	})

	t.Run("RecursiveRemove", func(t *testing.T) {

	})

	t.Run("MultiBlockRename", func(t *testing.T) {

	})

	t.Run("MultiBlockRemove", func(t *testing.T) {

	})

	t.Run("MultiBlockDirectoryWithFiles", func(t *testing.T) {

	})

	t.Run("MultiBlockRenameWithFiles", func(t *testing.T) {

	})

	t.Run("MultiBlockRemoveWithFiles", func(t *testing.T) {

	})
}

func createTestFS(t *testing.T, config *Config) (*LFS, tinyfs.BlockDevice, func()) {
	// create/format/mount the filesystem
	bd := tinyfs.NewMemoryDevice(testPageSize, testBlockSize, testBlockCount)
	fs := New(bd).Configure(config)
	if err := fs.Format(); err != nil {
		t.Error("Could not format", err)
	}
	if err := fs.Mount(); err != nil {
		t.Error("Could not mount", err)
	}
	return fs, bd, func() {
		if err := fs.Unmount(); err != nil {
			t.Error("Could not ummount", err)
		}
	}
}

func writeFileTest(t *testing.T, lfs *LFS, size int, name string) {
	buf := make([]byte, 32)
	f, err := lfs.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var totalBytes int
	for i := 0; i < size; i += len(buf) {
		b := buf[:]
		if len(b) > size-i {
			b = b[:size-i]
		}
		if _, err = rand.Read(b); err != nil {
			t.Fatal(err)
		}
		if n, err := f.Write(b); err != nil {
			t.Fatal(err)
		} else {
			totalBytes += n
		}
	}
	if totalBytes != size {
		t.Fatalf("expected to read %d bytes, was actually %d", size, totalBytes)
	}
}

func readFileTest(t *testing.T, lfs *LFS, size int, name string) {
	info, err := lfs.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != int64(size) {
		t.Fatalf("expected size %d but was actually %d", size, info.Size())
	}
	if info.IsDir() {
		t.Fatalf("expected file, but was a directory")
	}
	buf := make([]byte, 32)
	f, err := lfs.OpenFile(name, os.O_RDONLY)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var totalRead int
	for i := 0; i < size; i += len(buf) {
		b := buf[:]
		if len(b) > size-i {
			b = b[:size-i]
		}
		if n, err := f.Read(b); err != nil {
			t.Fatalf("error reading file after %d bytes: %v", n, err)
		} else {
			totalRead += n
		}
	}
	if totalRead != size {
		t.Fatalf("expected to read %d bytes, was actually %d", size, totalRead)
	}
	if _, err := f.Read(buf); err != io.EOF {
		t.Fatalf("expected io.EOF after %d bytes, was actually: %v", totalRead, err)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
