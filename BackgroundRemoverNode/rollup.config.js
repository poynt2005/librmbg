import typescript from "@rollup/plugin-typescript";
import commonjs from "@rollup/plugin-commonjs";
import terser from "@rollup/plugin-terser";

export default {
  input: "lib/js/index.ts",
  output: {
    file: "dist/BackgroundRemover.js",
    format: "cjs",
    sourcemap: false,
  },
  plugins: [
    typescript({
      tsconfig: "tsconfig.json",
    }),
    commonjs(),
    terser(),
  ],
};
