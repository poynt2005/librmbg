using System.Runtime.InteropServices;

namespace BackgroundRemover.Lib
{

    static class Utils
    {

        [DllImport("Urlmon.dll", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
        private static extern UInt64 FindMimeFromData(IntPtr pBC, string? pwzUrl, byte[] pBuffer, UInt64 cbSize, string? pwzMimeProposed, UInt64 dwMimeFlags, ref string ppwzMimeOut, UInt64 dwReserved);

        public static bool IsImage(this Byte[] b)
        {
            var mimeString = String.Empty;
            var ret = FindMimeFromData(IntPtr.Zero, null, b, (UInt64)b.Length, null, 0, ref mimeString, 0);

            if (ret != 0x0)
            {
                throw new Exception("cannot call FindMimeFromData from Urlmon.dll");
            }

            return mimeString.ToLower().StartsWith("image");
        }

    }


}
