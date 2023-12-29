import { checkRuntimeAndGetNativeBinding } from "./RuntimeChecker";

const nativeBinding: NativeBinding = checkRuntimeAndGetNativeBinding();

class BGRemoverCore {
  /**
   * This is a pointer to rembg instance in go context
   * @memberof BGRemoverCore
   */
  private m_strNativeHandle: string = "0";

  /**
   * Creates a BGRemoverCore instance, it will also create a handle point to go context
   * @constructor
   */
  constructor() {
    this.m_strNativeHandle = nativeBinding.RembgCreate();

    if (this.m_strNativeHandle == "0") {
      throw new Error("rembg instance creation encountered a unknown error");
    }
  }

  /**
   * Remember to call dispose when end of the operation,
   * This method will terminate and free the instance in go context
   * @memberof BGRemoverCore
   */
  public Dispose() {
    if (this.m_strNativeHandle == "0") {
      throw new Error("rembg instance has been disposed");
    }

    nativeBinding.RembgFree(this.m_strNativeHandle);
    this.m_strNativeHandle = "0";
  }

  /**
   * Remove the background from given image buffer
   * @param inputBuffer - image buffer you want to remove background with
   * @returns image buffer that is background removed
   */
  public RemoveBackground(inputBuffer: Buffer): Buffer {
    if (this.m_strNativeHandle == "0") {
      throw new Error("rembg instance not created");
    }

    return nativeBinding.RembgRemoveBackground(
      this.m_strNativeHandle,
      inputBuffer
    );
  }
}

export { BGRemoverCore };
