import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  build: {
    lib: {
      entry: 'src/index.js',     // 入口
      name: 'KefuChat',          // 全局变量名
      formats: ['umd'],          // 输出 UMD 格式（用于 <script>）
      fileName: () => 'kefu.min.js'
    },
    rollupOptions: {
      // 不 external Vue，要打包进去！
      external: [],
      output: {
        globals: {} // 不依赖外部 Vue
      }
    }
  },
})
