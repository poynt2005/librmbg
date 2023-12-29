#ifndef __RMBG_CONTEXT__DEF__H__
#define __RMBG_CONTEXT__DEF__H__

#include <stddef.h>
#include <stdint.h>

using pFnRembg_Create = uint64_t(*) Rembg_Create();
using pFnRembg_Free = void (*)(uint64_t cHandle);
using pFnRembg_RemoveBackground = int (*)(uint64_t cHandle, char *imgBufferInput, size_t imgBufferInputSize, char **imgBufferOutput, size_t *imgBufferOutputSize);
using pFnRembg_GetLastError = int (*)(uint64_t cHandle, char **outString);
using pFnRembg_ReleaseBuffer = void(*) Rembg_ReleaseBuffer(char **buffer);

using LibRmbgContext = struct libRmbgContext
{
    pFnRembg_Create Rembg_Create = nullptr;
    pFnRembg_Free Rembg_Free = nullptr;
    pFnRembg_RemoveBackground Rembg_RemoveBackground = nullptr;
    pFnRembg_GetLastError Rembg_GetLastError = nullptr;
    pFnRembg_ReleaseBuffer Rembg_ReleaseBuffer = nullptr;
};

enum class FunctionExecResult
{
    Success = 0,
    Fail,
    ContextNotCreated
};

#endif