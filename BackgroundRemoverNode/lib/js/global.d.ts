interface NativeBinding {
  RembgCreate: () => string;
  RembgFree: (handle: string) => null;
  RembgRemoveBackground: (handle: string, inputBuffer: Buffer) => Buffer;
  SystemDllCheck: (dllName: string) => boolean;
}
