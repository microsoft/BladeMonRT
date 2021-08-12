package test_utils

import "C"

func ToCString(guid string) *C.char {
	return C.CString(guid)
}