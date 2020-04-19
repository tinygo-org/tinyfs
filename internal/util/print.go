package util

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Fprintxxd(w io.Writer, offset uint32, b []byte) {
	var l int
	var buf16 = make([]byte, 16)
	var padding = ""
	for i, c := 0, len(b); i < c; i += 16 {
		l = i + 16
		if l >= c {
			padding = strings.Repeat(" ", (l-c)*3)
			l = c
		}
		_, _ = fmt.Fprintf(w, "%08x: % x    "+padding, offset+uint32(i), b[i:l])
		for j, n := 0, l-i; j < 16; j++ {
			if j >= n {
				buf16[j] = ' '
			} else if !strconv.IsPrint(rune(b[i+j])) {
				buf16[j] = '.'
			} else {
				buf16[j] = b[i+j]
			}
		}
		_, _ = w.Write(buf16)
		_, _ = fmt.Fprintln(w)
	}
}
