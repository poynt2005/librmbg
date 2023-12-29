import * as cp from "child_process";
import fs from "fs";
import { RuntimeError } from "./Errors";
import path from "path";
import { tmpdir } from "os";
import { uuidGen } from "./Utils";

const runtimeDll: string[] = [
  "BackgroundRemover.node",
  "librmbg.dll",
  "python3.dll",
  "python310.dll",
];

const systemRuntimeDll: string[] = ["vcruntime140.dll"];

const runtimeDllChecker = () => {
  runtimeDll.forEach((dll) => {
    if (!fs.existsSync(dll)) {
      throw new RuntimeError(`runtime dll: ${dll} is not exists`);
    }
  });
};

const systemButlinDllChecker = () => {
  const powershellCode: string = `
  $csCode = @"
  using System;
  using System.Runtime.InteropServices;
  using System.Collections.Generic;
  using System.Linq;
  namespace Utils {
    public static class Utils {
      [DllImport("kernel32", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
      private static extern IntPtr LoadLibraryA(string lpLibFileName);

      [DllImport("kernel32", CallingConvention = CallingConvention.Cdecl, CharSet = CharSet.Ansi)]
      private static extern IntPtr FreeLibrary(IntPtr hLibModule);

      public static bool CheckIsBuildinLibraryExists(string libName)
      {
        var hModule = LoadLibraryA(libName);
        var isExists = hModule != IntPtr.Zero;
        FreeLibrary(hModule);
        return isExists;
      }

      public static List<bool> CheckIsBuildinLibraryExists(List<string> libNames)
      {
        return libNames.Select(s => CheckIsBuildinLibraryExists(s)).ToList();
      }
    }
  }
"@
Add-Type -TypeDefinition $csCode -Language CSharp
$toChecks = [System.Collections.Generic.List[string]]@(${systemRuntimeDll.map(
    (dll) => `"${dll}"`
  )})
$result = [System.String]::Join(',', [Utils.Utils]::CheckIsBuildinLibraryExists($toChecks)) 
Write-Host $result
`;

  const tempPs: string = path.join(tmpdir(), uuidGen() + ".ps1");
  fs.writeFileSync(tempPs, powershellCode, "utf-8");
  const resultBuffer = cp.execSync(
    `powershell.exe -File "${tempPs}" -ExecutionPolicy RemoteSigned`
  );
  fs.rmSync(tempPs);

  const resultArr: boolean[] = JSON.parse(
    `[${new TextDecoder("utf-8").decode(resultBuffer).trim().toLowerCase()}]`
  ) as boolean[];

  systemRuntimeDll.forEach((dll, idx) => {
    if (!resultArr[idx]) {
      throw new RuntimeError(`butlin system runtime dll: ${dll} is not exists`);
    }
  });
};

export const checkRuntimeAndGetNativeBinding = (): NativeBinding => {
  runtimeDllChecker();
  systemButlinDllChecker();

  return require("./BackgroundRemover.node") as NativeBinding;
};
