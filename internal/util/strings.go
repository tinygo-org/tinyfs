package util

// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"unsafe"
)

// would be nice to use C.CString instead, but TinyGo doesn't seem to support
func CString(s string) unsafe.Pointer {
	ptr := C.malloc(C.size_t(len(s) + 1))
	buf := (*[1 << 28]byte)(ptr)[: len(s)+1 : len(s)+1]
	copy(buf, s)
	buf[len(s)] = 0
	return ptr
}

// would be nice to use C.GoString instead, but TinyGo doesn't seem to support
func GoString(s unsafe.Pointer) string {
	slen := int(C.strlen((*C.char)(s)))
	sbuf := make([]byte, slen)
	copy(sbuf, (*[1 << 28]byte)(unsafe.Pointer(s))[:slen:slen])
	return string(sbuf)
}
