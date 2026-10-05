// Bundles every src/tests/*.ts into dist/<name>.js for k6.
import { build } from "esbuild";
import { readdirSync } from "node:fs";

const entryPoints = readdirSync("src/tests")
  .filter((f) => f.endsWith(".ts"))
  .map((f) => `src/tests/${f}`);

await build({
  entryPoints,
  outdir: "dist",
  bundle: true,
  format: "esm",
  target: "es2022",
  platform: "neutral",
  external: ["k6", "k6/*"],
  logLevel: "info",
});
