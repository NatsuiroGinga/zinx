package util

import (
	"unsafe"
)

// String2Bytes converts string to []byte without copy.
func String2Bytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Bytes2String converts []byte to string without copy.
func Bytes2String(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}
