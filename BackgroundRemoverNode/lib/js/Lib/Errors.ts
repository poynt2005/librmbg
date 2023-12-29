export class RuntimeError extends Error {
  constructor(message: string) {
    super();
    this.name = "RuntimeError";
    this.message = message || "Error running the current script";
    this.stack = new Error().stack;
  }
}
