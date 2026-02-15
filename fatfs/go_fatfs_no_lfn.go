//go:build no_lfn

package fatfs

// FF_USE_LFN controls long file name (LFN) support in FatFs (default: 2).
// Using the build tag no_lfn sets FF_USE_LFN=0 which disables LFN and
// restricts file names to the 8.3 short file name format. It also removes
// the Unicode conversion tables in ffunicode.c, resulting in ~60KB of flash savings.

// #cgo CFLAGS: -DFF_USE_LFN=0
import "C"
