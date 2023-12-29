package main

import (
	"fmt"
	"librmbg/pycontext"
	"unsafe"
)

/*

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stddef.h>

enum FunctionExecResult {
	Success = 0,
	Fail,
	ContextNotCreated
};

*/
import "C"

type rembgContext struct {
	prr       *pycontext.PyRembgRemover
	errString string
}

var contextStore map[uint64]*rembgContext = map[uint64]*rembgContext{}

var globalLastErrorString string

//export Rembg_Create
func Rembg_Create() C.uint64_t {
	globalLastErrorString = ""

	prr, err := pycontext.NewPyRembgRemover()

	if err != nil {
		globalLastErrorString = err.Error()
		return 0
	}

	cHandle := uint64(uintptr(unsafe.Pointer(prr)))

	contextStore[cHandle] = &rembgContext{prr, ""}

	return C.uint64_t(cHandle)
}

//export Rembg_Free
func Rembg_Free(cHandle C.uint64_t) {
	globalLastErrorString = ""
	context, ok := contextStore[uint64(cHandle)]

	if ok {
		context.prr.Dispose()
		delete(contextStore, uint64(cHandle))
	}
}

//export Rembg_RemoveBackground
func Rembg_RemoveBackground(cHandle C.uint64_t, imgBufferInput *C.char, imgBufferInputSize C.size_t, imgBufferOutput **C.char, imgBufferOutputSize *C.size_t) C.int {
	globalLastErrorString = ""

	context, ok := contextStore[uint64(cHandle)]

	if !ok {
		globalLastErrorString = fmt.Sprintf("handle pointer: %d is not exists in context store", uint64(cHandle))
		return C.ContextNotCreated
	}
	context.errString = ""

	outputImage, err := context.prr.RemoveBackground(C.GoBytes(unsafe.Pointer(imgBufferInput), C.int(imgBufferInputSize)))

	if err != nil {
		context.errString = err.Error()
		return C.Fail
	}

	*imgBufferOutput = (*C.char)(C.CBytes(outputImage))
	*imgBufferOutputSize = (C.size_t)(len(outputImage))

	return C.Success
}

//export Rembg_GetLastError
func Rembg_GetLastError(cHandle C.uint64_t, outString **C.char) C.int {
	*outString = nil

	if uint64(cHandle) == 0 {
		*outString = C.CString(globalLastErrorString)
		return C.Success
	}

	context, ok := contextStore[uint64(cHandle)]

	if !ok {
		return C.ContextNotCreated
	}

	*outString = C.CString(context.errString)

	return C.Success
}

//export Rembg_ReleaseBuffer
func Rembg_ReleaseBuffer(buffer **C.char) {
	if buffer == nil || *buffer == nil {
		return
	}

	C.free(unsafe.Pointer(*buffer))
	*buffer = nil
}

func main() {}
