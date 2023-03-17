import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    disabled: process.env.NODE_ENV === "test"
  },
  plugins: [svelte()],
  publicDir: "public"
})
