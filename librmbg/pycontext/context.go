package pycontext

/*
#cgo CFLAGS: -I ../pylib/include
#cgo LDFLAGS: -L ../pylib/libs -lpython310

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stddef.h>
#include <stdint.h>
#include <Python.h>

void Py_DECREF_ALT(PyObject* o)
{
	Py_DECREF(o);
}

int PyBytes_Check_ALT(PyObject* o)
{
	return PyBytes_Check(o);
}

PyObject* CallRembgBgRemoveFunction(PyObject* pyRemoveFunction, PyObject* pyBytesInput)
{
	return PyObject_CallFunction(pyRemoveFunction, "O", pyBytesInput);
}

int PyBytes_AsStringAndSize_ALT(PyObject *obj, char **buffer, size_t *length)
{
	return PyBytes_AsStringAndSize(obj, buffer, length);
}

void weoij(const wchar_t* str)
{
	printf("%ls\n", str);
}
*/
import "C"
import (
	"embed"
	"fmt"
	"librmbg/utils"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

const _DLL_GUID = "5cb3360f-8bff-45a0-b334-ac988dccf566"

var pyRuntimePath string

//go:embed py_runtime.zip
var pyRuntimeZip embed.FS

var isPyInit bool = false

func pyInit() {
	if isPyInit {
		return
	}

	pyRuntimePath = filepath.Join(os.TempDir(), _DLL_GUID)
	pyHomePathStrW, err := windows.UTF16FromString(filepath.Join(pyRuntimePath, "python310.zip"))

	if err != nil {
		panic("cannot convert string to wstring")
	}

	pythonAppendScript := fmt.Sprintf("import sys;sys.path.insert(0, '');sys.path.append(r'%s');", pyRuntimePath)
	cStrPythonAppendScript := C.CString(pythonAppendScript)
	defer C.free(unsafe.Pointer(cStrPythonAppendScript))

	if info, err := os.Stat(pyRuntimePath); err != nil || !info.IsDir() {

		if os.MkdirAll(pyRuntimePath, os.ModePerm) != nil {
			panic("cannot make python runtime path")
		}

		fileContent, _ := pyRuntimeZip.ReadFile("py_runtime.zip")

		if utils.UnZipFromBuffer(fileContent, pyRuntimePath) != nil {
			panic("cannot unzip py runtime archive")
		}
	}

	C.Py_SetPythonHome((*C.wchar_t)(&pyHomePathStrW[0]))
	C.Py_SetPath((*C.wchar_t)(&pyHomePathStrW[0]))
	C.Py_Initialize()
	C.PyRun_SimpleString(cStrPythonAppendScript)

	isPyInit = true
}

type RembgModuleContext struct {
	rembg              *C.PyObject
	rembg_bg           *C.PyObject
	rembg_bg_remove_fn *C.PyObject
}

type PyRembgRemover struct {
	rembgContext *RembgModuleContext
}

func NewPyRembgRemover() (*PyRembgRemover, error) {
	pyInit()

	cStrRembgModuleName := C.CString("rembg")
	defer C.free(unsafe.Pointer(cStrRembgModuleName))

	pyRembgModule := C.PyImport_ImportModule(cStrRembgModuleName)

	if pyRembgModule == nil {
		return nil, fmt.Errorf("cannot import rembg module from python runtime")
	}

	cStrBackgroundAttrStr := C.CString("bg")
	defer C.free(unsafe.Pointer(cStrBackgroundAttrStr))

	pyRembgBackgroundModule := C.PyObject_GetAttrString(pyRembgModule, cStrBackgroundAttrStr)
	if pyRembgBackgroundModule == nil {
		C.Py_DECREF_ALT(pyRembgModule)
		return nil, fmt.Errorf("cannot import rembg.bg module from python runtime")
	}

	cStrRemoveAttrStr := C.CString("remove")
	defer C.free(unsafe.Pointer(cStrRemoveAttrStr))

	pyRembgBackgroundRemoveFunction := C.PyObject_GetAttrString(pyRembgBackgroundModule, cStrRemoveAttrStr)

	if pyRembgBackgroundRemoveFunction == nil {
		C.Py_DECREF_ALT(pyRembgBackgroundModule)

		return nil, fmt.Errorf("cannot import rembg.bg.remove function from python runtime")
	}

	if int(C.PyCallable_Check(pyRembgBackgroundRemoveFunction)) != 1 {
		C.Py_DECREF_ALT(pyRembgBackgroundRemoveFunction)
		C.Py_DECREF_ALT(pyRembgBackgroundModule)
		C.Py_DECREF_ALT(pyRembgModule)
		return nil, fmt.Errorf("rembg.bg.remove is not callable")
	}

	if !utils.ChkU2ModelDownload() {
		if err := utils.U2Download(); err != nil {
			return nil, fmt.Errorf("download u2 model failed, reason: %s", err.Error())
		}
	}

	return &PyRembgRemover{&RembgModuleContext{pyRembgModule, pyRembgBackgroundModule, pyRembgBackgroundRemoveFunction}}, nil
}

func (pc *PyRembgRemover) RemoveBackground(pictureBytes []byte) ([]byte, error) {
	if pc.rembgContext == nil {
		return nil, fmt.Errorf("rembg py context has been released")
	}

	if pictureBytes == nil || len(pictureBytes) == 0 {
		return nil, fmt.Errorf("cannt accept a empty byte")
	}

	cBytePictureBytes := (*C.char)(C.CBytes(pictureBytes))
	defer C.free(unsafe.Pointer(cBytePictureBytes))

	pyByteObject := C.PyBytes_FromStringAndSize(cBytePictureBytes, C.longlong(len(pictureBytes)))

	if pyByteObject == nil {
		return nil, fmt.Errorf("cannot marshal from bytes to PyBytesObject")
	}
	defer C.Py_DECREF_ALT(pyByteObject)

	cStrArgumentFormat := C.CString("O")
	defer C.free(unsafe.Pointer(cStrArgumentFormat))

	pyBackgroundResultObject := C.CallRembgBgRemoveFunction(pc.rembgContext.rembg_bg_remove_fn, pyByteObject)

	if pyBackgroundResultObject == nil {
		return nil, fmt.Errorf("cannot call rembg.bg.remove function")
	}
	defer C.Py_DECREF_ALT(pyBackgroundResultObject)

	if int(C.PyBytes_Check_ALT(pyBackgroundResultObject)) == 0 {
		return nil, fmt.Errorf("rembg.bg.remove returned a invalid result")
	}

	var pyBackgroundResultBuffer *C.char = nil
	var cResultBufferSize C.size_t = 0

	if int(C.PyBytes_AsStringAndSize_ALT(pyBackgroundResultObject, (**C.char)(&pyBackgroundResultBuffer), (*C.size_t)(&cResultBufferSize))) != 0 {
		return nil, fmt.Errorf("cannot marshal from py bytes object to c buffer")
	}

	if pyBackgroundResultBuffer == nil || cResultBufferSize == 0 {
		return nil, fmt.Errorf("py bytes returned a empty buffer")
	}

	return C.GoBytes(unsafe.Pointer(pyBackgroundResultBuffer), C.int(cResultBufferSize)), nil
}

func (pc *PyRembgRemover) Dispose() {
	if pc.rembgContext == nil {
		return
	}

	if pc.rembgContext.rembg_bg_remove_fn != nil {
		C.Py_DECREF_ALT(pc.rembgContext.rembg_bg_remove_fn)
		pc.rembgContext.rembg_bg_remove_fn = nil
	}

	if pc.rembgContext.rembg_bg != nil {
		C.Py_DECREF_ALT(pc.rembgContext.rembg_bg)
		pc.rembgContext.rembg_bg = nil
	}

	if pc.rembgContext.rembg != nil {
		C.Py_DECREF_ALT(pc.rembgContext.rembg)
		pc.rembgContext.rembg = nil
	}

	pc.rembgContext = nil
}
