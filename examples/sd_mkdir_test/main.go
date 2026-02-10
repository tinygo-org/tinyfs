// SD Card LittleFS Console — interactive serial console for LittleFS on SD card.
//
// This is a standalone example (does not use the shared console package) that
// demonstrates LittleFS on an SD card over SPI. It provides basic filesystem
// commands over the USB serial port at 115200 baud.
//
// Build for Grand Central M4:
//
//	tinygo build -target=grandcentral-m4 -stack-size=16KB -o firmware.uf2 .
//
// Commands: HELP, FORMAT, MOUNT, UNMOUNT, LS, MKDIR, WRITE, CAT, RM,
//
//	TEST (basic mkdir), TEST2 (nested dirs), TEST3 (files), TEST4 (stress)
package main

import (
	"fmt"
	"machine"
	"os"
	"strings"
	"time"

	"tinygo.org/x/drivers/sdcard"
	"tinygo.org/x/tinyfs/littlefs"
)

var (
	sd      *sdcard.Device
	fs      *littlefs.LFS
	mounted bool
)

func main() {
	time.Sleep(2 * time.Second)
	fmt.Println("=== SD Card LittleFS Console ===")
	fmt.Println("Type HELP for commands")

	// Init SPI + SD card
	machine.SPI1.Configure(machine.SPIConfig{
		SCK: machine.SDCARD_SCK_PIN, SDO: machine.SDCARD_SDO_PIN,
		SDI: machine.SDCARD_SDI_PIN, Frequency: 1000000,
	})
	machine.SDCARD_CS_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SDCARD_CS_PIN.High()

	dev := sdcard.New(machine.SPI1,
		machine.SDCARD_SCK_PIN, machine.SDCARD_SDO_PIN,
		machine.SDCARD_SDI_PIN, machine.SDCARD_CS_PIN)
	sd = &dev
	if err := sd.Configure(); err != nil {
		fmt.Println("SD init failed:", err)
		select {}
	}

	fs = littlefs.New(sd)
	fs.Configure(&littlefs.Config{
		CacheSize: 512, LookaheadSize: 512, BlockCycles: 100,
	})
	if err := fs.Mount(); err != nil {
		fmt.Println("Mount failed (try FORMAT):", err)
	} else {
		mounted = true
		fmt.Println("Mounted OK")
	}

	// Command loop
	var buf strings.Builder
	for {
		if machine.Serial.Buffered() > 0 {
			c, _ := machine.Serial.ReadByte()
			if c == '\n' || c == '\r' {
				if cmd := strings.TrimSpace(buf.String()); cmd != "" {
					run(cmd)
				}
				buf.Reset()
			} else {
				buf.WriteByte(c)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func run(cmd string) {
	parts := strings.SplitN(cmd, " ", 2)
	action := strings.ToUpper(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}
	fmt.Println(">>>", cmd)

	switch action {
	case "HELP":
		fmt.Println("  FORMAT  MOUNT  UNMOUNT  LS [path]  MKDIR path")
		fmt.Println("  WRITE path data  CAT path  RM path")
		fmt.Println("  TEST   — 10x mkdir+stat (basic bug trigger)")
		fmt.Println("  TEST2  — nested directories (parent/child/grandchild)")
		fmt.Println("  TEST3  — files inside directories")
		fmt.Println("  TEST4  — 50x high-frequency mkdir+stat stress")
	case "FORMAT":
		if mounted {
			fs.Unmount()
			mounted = false
		}
		if err := fs.Format(); err != nil {
			fmt.Println("Format FAILED:", err)
			return
		}
		if err := fs.Mount(); err != nil {
			fmt.Println("Mount FAILED:", err)
			return
		}
		mounted = true
		fmt.Println("OK")
	case "MOUNT":
		if mounted {
			fmt.Println("already mounted")
			return
		}
		if err := fs.Mount(); err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		mounted = true
		fmt.Println("OK")
	case "UNMOUNT":
		if !mounted {
			fmt.Println("not mounted")
			return
		}
		if err := fs.Unmount(); err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		mounted = false
		fmt.Println("OK")
	case "LS":
		if !requireMount() {
			return
		}
		if arg == "" {
			arg = "/"
		}
		dir, err := fs.Open(arg)
		if err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		entries, _ := dir.Readdir(-1)
		dir.Close()
		for _, e := range entries {
			t := "f"
			if e.IsDir() {
				t = "d"
			}
			fmt.Printf("  %s %8d %s\n", t, e.Size(), e.Name())
		}
		if len(entries) == 0 {
			fmt.Println("  (empty)")
		}
	case "MKDIR":
		if !requireMount() || requireArg(arg) {
			return
		}
		if err := fs.Mkdir(arg, 0755); err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		fmt.Println("OK")
	case "WRITE":
		if !requireMount() {
			return
		}
		wp := strings.SplitN(arg, " ", 2)
		if len(wp) < 2 {
			fmt.Println("Usage: WRITE path data")
			return
		}
		f, err := fs.OpenFile(wp[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
		if err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		f.Write([]byte(wp[1]))
		f.Close()
		fmt.Println("OK")
	case "CAT":
		if !requireMount() || requireArg(arg) {
			return
		}
		f, err := fs.Open(arg)
		if err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		buf := make([]byte, 256)
		n, _ := f.Read(buf)
		f.Close()
		fmt.Printf("%s\n", buf[:n])
	case "RM":
		if !requireMount() || requireArg(arg) {
			return
		}
		if err := fs.Remove(arg); err != nil {
			fmt.Println("FAILED:", err)
			return
		}
		fmt.Println("OK")
	case "TEST":
		testMkdir()
	case "TEST2":
		testNested()
	case "TEST3":
		testFiles()
	case "TEST4":
		testStress()
	default:
		fmt.Println("Unknown command. Type HELP.")
	}
}

func requireMount() bool {
	if !mounted {
		fmt.Println("not mounted")
	}
	return mounted
}

func requireArg(arg string) bool {
	if arg == "" {
		fmt.Println("missing argument")
		return true
	}
	return false
}

// testMkdir exercises mkdir+stat in a loop — the pattern that triggers the
// LLVM -O2 miscompilation bug in lfs_dir_fetchmatch (see littlefs/lfs.c).
// At -O2 without the optnone fix, every stat fails with "no directory entry".
func testMkdir() {
	if !requireMount() {
		return
	}
	const n = 10
	fail := 0
	for i := 0; i < n; i++ {
		p := fmt.Sprintf("/t_%d", i)
		if err := fs.Mkdir(p, 0755); err != nil {
			fmt.Printf("  mkdir %s: %v\n", p, err)
			fail++
			continue
		}
		if info, err := fs.Stat(p); err != nil {
			fmt.Printf("  FAIL %s: stat: %v\n", p, err)
			fail++
		} else if !info.IsDir() {
			fmt.Printf("  FAIL %s: not a dir\n", p)
			fail++
		} else {
			fmt.Printf("  OK   %s\n", p)
		}
	}
	for i := 0; i < n; i++ {
		fs.Remove(fmt.Sprintf("/t_%d", i))
	}
	if fail > 0 {
		fmt.Printf("RESULT: %d/%d FAILED\n", fail, n)
	} else {
		fmt.Printf("RESULT: %d/%d passed\n", n, n)
	}
}

// testNested creates parent → child → grandchild directories and verifies
// each level with stat, then cleans up in reverse order.
func testNested() {
	if !requireMount() {
		return
	}
	levels := []string{"/na", "/na/nb", "/na/nb/nc"}
	fail := 0
	for _, p := range levels {
		if err := fs.Mkdir(p, 0755); err != nil {
			fmt.Printf("  mkdir %s: %v\n", p, err)
			fail++
			continue
		}
		if info, err := fs.Stat(p); err != nil {
			fmt.Printf("  FAIL %s: stat: %v\n", p, err)
			fail++
		} else if !info.IsDir() {
			fmt.Printf("  FAIL %s: not a dir\n", p)
			fail++
		} else {
			fmt.Printf("  OK   %s\n", p)
		}
	}
	// cleanup reverse
	for i := len(levels) - 1; i >= 0; i-- {
		fs.Remove(levels[i])
	}
	if fail > 0 {
		fmt.Printf("RESULT: %d/%d FAILED\n", fail, len(levels))
	} else {
		fmt.Printf("RESULT: %d/%d passed\n", len(levels), len(levels))
	}
}

// testFiles creates directories and writes+reads files inside them.
func testFiles() {
	if !requireMount() {
		return
	}
	fail := 0
	dirs := []string{"/fd1", "/fd1/fd2"}
	for _, d := range dirs {
		if err := fs.Mkdir(d, 0755); err != nil {
			fmt.Printf("  mkdir %s: %v\n", d, err)
			fail++
		}
	}
	// write files inside dirs
	files := map[string]string{
		"/fd1/a.txt":     "hello",
		"/fd1/fd2/b.txt": "world",
	}
	for path, data := range files {
		f, err := fs.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
		if err != nil {
			fmt.Printf("  FAIL write %s: %v\n", path, err)
			fail++
			continue
		}
		f.Write([]byte(data))
		f.Close()
		// read back
		rf, err := fs.Open(path)
		if err != nil {
			fmt.Printf("  FAIL open %s: %v\n", path, err)
			fail++
			continue
		}
		buf := make([]byte, 64)
		n, _ := rf.Read(buf)
		rf.Close()
		if string(buf[:n]) != data {
			fmt.Printf("  FAIL %s: got %q want %q\n", path, string(buf[:n]), data)
			fail++
		} else {
			fmt.Printf("  OK   %s\n", path)
		}
	}
	// cleanup
	for path := range files {
		fs.Remove(path)
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		fs.Remove(dirs[i])
	}
	total := len(files)
	if fail > 0 {
		fmt.Printf("RESULT: %d/%d FAILED\n", fail, total)
	} else {
		fmt.Printf("RESULT: %d/%d passed\n", total, total)
	}
}

// testStress does 50 rapid mkdir+stat cycles to stress-test the fix.
func testStress() {
	if !requireMount() {
		return
	}
	const n = 50
	fail := 0
	for i := 0; i < n; i++ {
		p := fmt.Sprintf("/s_%d", i)
		if err := fs.Mkdir(p, 0755); err != nil {
			fmt.Printf("  mkdir %s: %v\n", p, err)
			fail++
			continue
		}
		if _, err := fs.Stat(p); err != nil {
			fmt.Printf("  FAIL %s: %v\n", p, err)
			fail++
		}
	}
	for i := 0; i < n; i++ {
		fs.Remove(fmt.Sprintf("/s_%d", i))
	}
	if fail > 0 {
		fmt.Printf("RESULT: %d/%d FAILED\n", fail, n)
	} else {
		fmt.Printf("RESULT: %d/%d passed\n", n, n)
	}
}
