using System.Runtime.InteropServices;

namespace BackgroundRemover
{
    public enum FunctionExecResult
    {
        Success = 0,
        Fail,
        ContextNotCreated
    };

    public class BGRemover : IDisposable
    {
        [DllImport("librmbg", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        private static extern UInt64 Rembg_Create();

        [DllImport("librmbg", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        private static extern void Rembg_Free(UInt64 cHandle);
        [DllImport("librmbg", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        private static extern int Rembg_RemoveBackground(UInt64 cHandle, IntPtr imgBufferInput, UInt64 imgBufferInputSize, ref IntPtr imgBufferOutput, ref UInt64 imgBufferOutputSize);
        [DllImport("librmbg", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        private static extern int Rembg_GetLastError(UInt64 cHandle, ref IntPtr outString);
        [DllImport("librmbg", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
        private static extern void Rembg_ReleaseBuffer(ref IntPtr buffer);

        private UInt64 GoHandle { get; set; } = 0;
        private bool Disposed { get; set; } = false;

        public BGRemover()
        {
            GoHandle = Rembg_Create();

            if (GoHandle == 0)
            {
                var errReasonPtr = IntPtr.Zero;
                var errRet = (FunctionExecResult)Rembg_GetLastError(0, ref errReasonPtr);

                if (errReasonPtr == IntPtr.Zero)
                {
                    throw new Exception("construct a rembg context encountered a unknown error");
                }

                var errReason = Marshal.PtrToStringAnsi(errReasonPtr);
                Rembg_ReleaseBuffer(ref errReasonPtr);

                throw new Exception(errReason);
            }
        }

        public byte[] RemoveBackground(string sourceImage)
        {
            return RemoveBackground(File.ReadAllBytes(sourceImage));
        }

        public void RemoveBackground(string sourceImage, string destImage)
        {
            var destImageBuffer = RemoveBackground(File.ReadAllBytes(sourceImage));
            File.WriteAllBytes(destImage, destImageBuffer);
        }

        public byte[] RemoveBackground(byte[] sourceImage)
        {
            if (sourceImage.Length == 0)
            {
                throw new Exception("cannot remove  background with a empty image content");
            }

            var unmanagedSourceImageBuffer = Marshal.AllocHGlobal(sourceImage.Length);
            Marshal.Copy(sourceImage, 0, unmanagedSourceImageBuffer, sourceImage.Length);

            var unmanagedDestImageBuffer = IntPtr.Zero;
            UInt64 destImageBufferSize = 0;

            InterpretException((FunctionExecResult)Rembg_RemoveBackground(GoHandle, unmanagedSourceImageBuffer, (UInt64)sourceImage.Length, ref unmanagedDestImageBuffer, ref destImageBufferSize));

            if (unmanagedDestImageBuffer == IntPtr.Zero || destImageBufferSize == 0)
            {
                throw new Exception("cannot get any buffer from go context");
            }

            var destImageBuffer = new byte[destImageBufferSize];

            Marshal.Copy(unmanagedDestImageBuffer, destImageBuffer, 0, (int)destImageBufferSize);
            Marshal.FreeHGlobal(unmanagedSourceImageBuffer);

            return destImageBuffer;
        }

        private void InterpretException(FunctionExecResult ret)
        {
            if (ret == FunctionExecResult.ContextNotCreated)
            {
                throw new Exception($"rembg instance is not created in go context");
            }
            else if (ret == FunctionExecResult.Fail)
            {
                var errReasonPtr = IntPtr.Zero;
                var errReason = $"rembg function call returned a unknown reason error";

                var errRet = (FunctionExecResult)Rembg_GetLastError(GoHandle, ref errReasonPtr);

                if (errReasonPtr != IntPtr.Zero)
                {
                    errReason = Marshal.PtrToStringAnsi(errReasonPtr);
                    Rembg_ReleaseBuffer(ref errReasonPtr);
                }

                throw new Exception(errReason);
            }
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!Disposed)
            {
                if (disposing)
                {
                    //dispose managed resources
                    if (GoHandle != 0)
                    {
                        Rembg_Free(GoHandle);
                        GoHandle = 0;
                    }
                }
            }
            //dispose unmanaged resources
            Disposed = true;
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }
    }


}