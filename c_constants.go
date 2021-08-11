package main

import "C"

func GUIDForTest() *C.char {
	return C.CString("50bd065e-f3e9-4887-8093-b171f1b01372")
}