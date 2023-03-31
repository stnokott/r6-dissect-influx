import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { warmup } from 'vite-plugin-warmup'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
    warmup({
      clientFiles: ['./cypress/support/component.ts', './src/**/*.cy.{js,jsx,ts,tsx}']
    })
  ],
  publicDir: "public"
})
