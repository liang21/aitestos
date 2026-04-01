import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import arcoReactPlugin from '@arco-plugins/vite-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    arcoReactPlugin({
      style: 'css',
    }),
  ],
})
