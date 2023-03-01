import sveltePreprocess from 'svelte-preprocess'
import { optimizeCss, optimizeImports } from 'carbon-preprocess-svelte'

import { dirname, join } from 'path'
import { fileURLToPath } from 'url'
const __dirname = dirname(fileURLToPath(import.meta.url))

export default {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: [
    sveltePreprocess({
      scss: {
        includePaths: [
          join(__dirname, 'node_modules')
        ]
      }
    }),
    optimizeImports(),
    optimizeCss(),
  ]
}
