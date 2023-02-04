import sveltePreprocess from 'svelte-preprocess'
import { optimizeCss, optimizeImports } from 'carbon-preprocess-svelte'

export default {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: [
    sveltePreprocess(),
    optimizeImports(),
    optimizeCss(),
  ]
}
