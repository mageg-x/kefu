import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

function copyAssetsPlugin() {
  return {
    name: 'copy-assets',
    closeBundle() {
      const source = path.resolve(__dirname, 'src/assets/logo.png')
      const target = path.resolve(__dirname, 'dist/logo.png')
      if (fs.existsSync(source)) {
        fs.copyFileSync(source, target)
      }
    }
  }
}

function copyChatPlugin() {
  return {
    name: 'copy-chat-html',
    closeBundle() {
      const source = path.resolve(__dirname, 'chat.html')
      const target = path.resolve(__dirname, 'dist/chat.html')
      if (fs.existsSync(source)) {
        fs.copyFileSync(source, target)
      }
    }
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    copyAssetsPlugin(),
    copyChatPlugin()
  ],
  define: {
    'process.env.NODE_ENV': JSON.stringify('production'),
    'process.env': '{}',
  },
  publicDir: false,
  build: {
    outDir: 'dist',
    cssCodeSplit: false,
    rollupOptions: {
      input: {
        widget: 'src/widget.js',
        chat: 'src/chat.js'
      },
      output: {
        entryFileNames: chunk => {
            if (chunk.name === 'widget') {
              return 'widget.min.js'
            } else if (chunk.name === 'chat') {
              return 'chat.min.js'
            }
            return '[name].[hash].js'
          },
        chunkFileNames: 'vendor.js',
        assetFileNames: assetInfo => {
          if (assetInfo.name === 'style.css') {
            return 'style.css'
          }
          return '[name].[ext]'
        },
        globals: {}
      }
    }
  },
})
