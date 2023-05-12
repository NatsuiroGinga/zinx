package convert

import "unsafe"

// ToBytes converts string to []byte without copy.
func ToBytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ToString converts []byte to string without copy.
func ToString(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}
