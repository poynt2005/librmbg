#include <Napi.h>
#include <Windows.h>

#include <optional>
#include <string>
#include <memory>

#include "RmbgContextDef.h"

static LibRmbgContext libRmbgContext;

#define LOAD_CONTEXT(MODULE, CONTEXT, FN)                               \
    CONTEXT.##FN = (pFn##FN)GetProcAddress(MODULE, #FN);                \
    if (CONTEXT.##FN == nullptr)                                        \
    {                                                                   \
        return "cannot load function \"" #FN "\" from dynamic library"; \
    }

#define CHECK_ARGUMENTS_MIN_COUNT(COUNT)                                                                                          \
    if (info.Length() < COUNT)                                                                                                    \
    {                                                                                                                             \
        Napi::TypeError::New(env, "invalid count of arguments, must at least " #COUNT " arguments").ThrowAsJavaScriptException(); \
        return env.Null();                                                                                                        \
    }

#define CHECK_ARGUMENTS_TYPE(TYPE, ASTYPE, POS)                                                                         \
    if (!info[POS].Is##TYPE())                                                                                          \
    {                                                                                                                   \
        Napi::TypeError::New(env, "parameter position " #POS " must be a " #TYPE " type").ThrowAsJavaScriptException(); \
        return env.Null();                                                                                              \
    }                                                                                                                   \
    auto arg##POS = info[POS].As<Napi::##ASTYPE>();

#define CONVERT_ARG0_TO_UHANDLE \
    uint64_t uHandle = std::stoull(arg0.Utf8Value());

#define NODE_FUNC(FUNCNAME) \
    exports.Set(Napi::String::New(env, #FUNCNAME), Napi::Function::New(env, NodeFunc_##FUNCNAME));

std::optional<std::string> LoadDllContext()
{
    auto rmbgModule = LoadLibrary("librmbg.dll");

    if (rmbgModule == nullptr)
    {
        return "cannot load dynamic library \"librmbg.dll\"";
    }

    LOAD_CONTEXT(rmbgModule, libRmbgContext, Rembg_Create);
    LOAD_CONTEXT(rmbgModule, libRmbgContext, Rembg_Free);
    LOAD_CONTEXT(rmbgModule, libRmbgContext, Rembg_RemoveBackground);
    LOAD_CONTEXT(rmbgModule, libRmbgContext, Rembg_GetLastError);
    LOAD_CONTEXT(rmbgModule, libRmbgContext, Rembg_ReleaseBuffer);

    return std::nullopt;
}
Napi::Value NodeFunc_RembgCreate(const Napi::CallbackInfo &info)
{
    auto env = info.Env();

    auto uHandle = libRmbgContext.Rembg_Create();

    if (!uHandle)
    {
        std::string strErrorReason("rembg creation encountered a unknown error");

        char *szErrorReason = nullptr;
        libRmbgContext.Rembg_GetLastError(0, &szErrorReason);

        if (szErrorReason != nullptr)
        {
            strErrorReason = szErrorReason;
            libRmbgContext.Rembg_ReleaseBuffer(&szErrorReason);
        }

        Napi::Error::New(env, strErrorReason).ThrowAsJavaScriptException();
        return env.Null();
    }

    return Napi::String::New(env, std::to_string(uHandle));
}

Napi::Value NodeFunc_RembgFree(const Napi::CallbackInfo &info)
{
    auto env = info.Env();
    CHECK_ARGUMENTS_MIN_COUNT(1);
    CHECK_ARGUMENTS_TYPE(String, String, 0);
    CONVERT_ARG0_TO_UHANDLE;
    libRmbgContext.Rembg_Free(uHandle);
    return env.Null();
}

Napi::Value NodeFunc_RembgRemoveBackground(const Napi::CallbackInfo &info)
{
    auto env = info.Env();
    CHECK_ARGUMENTS_MIN_COUNT(2);
    CHECK_ARGUMENTS_TYPE(String, String, 0);
    CHECK_ARGUMENTS_TYPE(Buffer, Buffer<char>, 1);
    CONVERT_ARG0_TO_UHANDLE;

    char *pDestBuffer = nullptr;
    size_t uDestBufferSize = 0;

    auto ret = (FunctionExecResult)libRmbgContext.Rembg_RemoveBackground(uHandle, arg1.Data(), arg1.Length(), &pDestBuffer, &uDestBufferSize);

    if (ret != FunctionExecResult::Success)
    {
        if (ret == FunctionExecResult::ContextNotCreated)
        {
            Napi::Error::New(env, "rambg context is not created").ThrowAsJavaScriptException();
            return env.Null();
        }
        else
        {
            std::string strErrorReason("calling RemoveBackground function encountered a unknown error");

            char *szErrorReason = nullptr;
            libRmbgContext.Rembg_GetLastError(uHandle, &szErrorReason);

            if (szErrorReason != nullptr)
            {
                strErrorReason = szErrorReason;
                libRmbgContext.Rembg_ReleaseBuffer(&szErrorReason);
            }

            Napi::Error::New(env, strErrorReason).ThrowAsJavaScriptException();
            return env.Null();
        }
    }

    if (pDestBuffer == nullptr || !uDestBufferSize)
    {
        Napi::Error::New(env, "node context did not receieve ANY file content from go context").ThrowAsJavaScriptException();
        return env.Null();
    }

    auto kDestBuffer = Napi::Buffer<char>::Copy(env, pDestBuffer, uDestBufferSize);
    libRmbgContext.Rembg_ReleaseBuffer(&pDestBuffer);

    return kDestBuffer;
}

Napi::Object Initialize(Napi::Env env, Napi::Object exports)
{
    auto loadErr = LoadDllContext();

    if (loadErr.has_value())
    {
        Napi::Error::New(env, loadErr.value()).ThrowAsJavaScriptException();
        return Napi::Object::New(env);
    }

    NODE_FUNC(RembgCreate);
    NODE_FUNC(RembgFree);
    NODE_FUNC(RembgRemoveBackground);

    return exports;
}

NODE_API_MODULE(BackgroundRemover, Initialize)
