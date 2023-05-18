# TinyFS 

[![Build](https://github.com/tinygo-org/tinyfs/actions/workflows/build.yml/badge.svg?branch=dev)](https://github.com/tinygo-org/tinyfs/actions/workflows/build.yml)

TinyFS contains Go implementations of embedded filesystems. The packages in
this module require CGo support and are [TinyGo](https://tinygo.org/) compatible.

## Supported hardware

You can use TinyFS with the following embedded hardware configurations:

### Onboard Flash memory

You can use TinyFS on processors that have onboard Flash memory in the area above where the compiled TinyGo program is running.

### External Flash memory

You can use TinyFS on an external Flash memory chip connected via SPI/QSPI hardware interface.

See https://github.com/tinygo-org/drivers/tree/release/flash

### SD card connected via SPI/QPI

You can use TinyFS on an SD card connected via SPI/QSPI hardware interface.

See https://github.com/tinygo-org/drivers/tree/release/sdcard

## LittleFS

The LittleFS file system is specifically designed for embedded applications.

See https://github.com/littlefs-project/littlefs for more information.

### Example

This example runs on the RP2040 using the on-board flash in the available memory above where the program code itself is running:

```
$ tinygo flash -target pico -monitor ./examples/console/littlefs/machine/
Connected to /dev/ttyACM0. Press Ctrl-C to exit.
SPI Configured. Reading flash info
==> lsblk

-------------------------------------
 Device Information:  
-------------------------------------
 flash data start: 10017000
 flash data end:  10200000
-------------------------------------

==> format
Successfully formatted LittleFS filesystem.

==> mount
Successfully mounted LittleFS filesystem.

==> ls
==> samples
wrote 90 bytes to file0.txt
wrote 90 bytes to file1.txt
wrote 90 bytes to file2.txt
wrote 90 bytes to file3.txt
wrote 90 bytes to file4.txt
==> ls
-rwxrwxrwx    90 file0.txt
-rwxrwxrwx    90 file1.txt
-rwxrwxrwx    90 file2.txt
-rwxrwxrwx    90 file3.txt
-rwxrwxrwx    90 file4.txt
==> cat file3.txt
00000000: 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f     !"#$%&'()*+,-./
00000010: 30 31 32 33 34 35 36 37 38 39 3a 3b 3c 3d 3e 3f    0123456789:;<=>?
00000020: 40 41 42 43 44 45 46 47 48 49 4a 4b 4c 4d 4e 4f    @ABCDEFGHIJKLMNO
00000030: 50 51 52 53 54 55 56 57 58 59 5a 5b 5c 5d 5e 5f    PQRSTUVWXYZ[\]^_
00000040: 60 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f    `abcdefghijklmno
00000050: 70 71 72 73 74 75 76 77 78 79                      pqrstuvwxy
```

After unplugging and reconnecting the RP2040 device (a hard restart):

```
$ tinygo monitor
Connected to /dev/ttyACM0. Press Ctrl-C to exit.
SPI Configured. Reading flash info
==> mount
Successfully mounted LittleFS filesystem.

==> ls
-rwxrwxrwx    90 file0.txt
-rwxrwxrwx    90 file1.txt
-rwxrwxrwx    90 file2.txt
-rwxrwxrwx    90 file3.txt
-rwxrwxrwx    90 file4.txt
==> cat file3.txt
00000000: 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f     !"#$%&'()*+,-./
00000010: 30 31 32 33 34 35 36 37 38 39 3a 3b 3c 3d 3e 3f    0123456789:;<=>?
00000020: 40 41 42 43 44 45 46 47 48 49 4a 4b 4c 4d 4e 4f    @ABCDEFGHIJKLMNO
00000030: 50 51 52 53 54 55 56 57 58 59 5a 5b 5c 5d 5e 5f    PQRSTUVWXYZ[\]^_
00000040: 60 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f    `abcdefghijklmno
00000050: 70 71 72 73 74 75 76 77 78 79                      pqrstuvwxy
```

## FAT FS

The FAT file system is not currently working, due to https://github.com/tinygo-org/tinygo/issues/3460.

