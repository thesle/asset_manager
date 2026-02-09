import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webStoresPath = path.resolve(__dirname, '../../web/src/stores.js');
const desktopStoresPath = path.resolve(__dirname, 'src/stores.js');

// Plugin to redirect web stores imports to desktop stores
function redirectWebStores() {
  return {
    name: 'redirect-web-stores',
    enforce: 'pre',
    resolveId(source, importer) {
      if (importer && source.endsWith('stores.js')) {
        const resolved = path.resolve(path.dirname(importer), source);
        if (resolved === webStoresPath) {
          return desktopStoresPath;
        }
      }
      return null;
    }
  };
}

export default defineConfig({
  plugins: [redirectWebStores(), svelte()],
  resolve: {
    alias: {
      '@shared': path.resolve(__dirname, '../../shared'),
    },
  },
});
