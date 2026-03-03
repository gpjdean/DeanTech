import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      // 只代理API请求，以/api开头
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        // 不要重写路径，保留/api前缀
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})