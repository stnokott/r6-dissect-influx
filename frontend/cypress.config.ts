import { defineConfig } from "cypress";
import watch from "@cypress/watch-preprocessor";

export default defineConfig({
  component: {
    devServer: {
      framework: "svelte",
      bundler: "vite",
    },
    setupNodeEvents(on, _config) {
      on("file:preprocessor", watch())
    },
  },
  video: false
});
