import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
  return {
    plugins: [react()],
    server: {
      proxy:
        mode === 'development'
          ? {
              '/api': {
                target: 'http://localhost:5001',
                changeOrigin: true,
              },
            }
          : undefined,
    },
  }
})
