import { build } from "bun"
import { watch } from "node:fs"

console.log("Starting file watcher at:", new URL("src", import.meta.url).pathname)

const watcher = watch(
  new URL("src", import.meta.url).pathname,
  { recursive: true },
  async (event, filename) => {
    console.log(`Detected ${event} in ${filename}`)
    await build({
      entrypoints: [new URL(`src/${filename}`, import.meta.url).pathname],
      outdir: "./build/js",
      root: "./src",
      target: "node",
    })
  }
)

process.on("SIGINT", () => {
  // close watcher when Ctrl-C is pressed
  console.log("Closing watcher...")
  watcher.close()

  process.exit(0)
})
