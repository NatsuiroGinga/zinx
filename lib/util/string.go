package util

import (
	"unsafe"
)

func String2Bytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Bytes2String(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}
