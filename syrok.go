package main

/*
#cgo LDFLAGS: -lsyrok -lm
#include <syrok.h>
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

func syrokImage(data []byte, mode int) ([]byte, error) {
	var size C.int
	var err *C.uchar
	var data_ []byte

	err = C.syrok((*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(len(data)), &size, C.int(mode))
	if unsafe.Pointer(err) == C.NULL {
		return nil, errors.New(C.GoString(C.syrok_get_error()))
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data_))
	sh.Data = uintptr(unsafe.Pointer(err))
	sh.Len = int(size)
	sh.Cap = int(size)
	data_ = *(*[]byte)(unsafe.Pointer(sh))
	return data_, nil
}
