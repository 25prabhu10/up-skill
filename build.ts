import { build, Glob } from "bun"

const glob = new Glob("./src/**/*.ts")

const entrypoints = await Array.fromAsync(glob.scan("."))

await build({
  entrypoints,
  outdir: "./build/js",
  root: "./src",
  target: "node",
})
