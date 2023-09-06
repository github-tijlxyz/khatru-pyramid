import { defineConfig } from "vite";
import { svelte, vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import postcss from "postcss";

// see https://vitejs.dev/config
export default defineConfig({
  envPrefix: "UI",
  base: "./",
  build: {
    chunkSizeWarningLimit: 1000,
    reportCompressedSize: false,
  },
  plugins: [
    svelte({
      preprocess: [vitePreprocess()],
    }),
  ],
  css: {
    postcss,
  },
  resolve: {
    alias: {
      "@": __dirname + "/src",
    },
  },
});
