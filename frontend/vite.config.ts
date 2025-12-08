import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { configDefaults } from 'vitest/config';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  
  return {
    plugins: [react()],
    
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },

    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: ['./src/test/setup.ts'],
      css: true,
      coverage: {
        provider: 'v8',
        reporter: ['text', 'json', 'html'],
        exclude: [
          ...configDefaults.coverage.exclude,
          'src/test/**',
          '**/*.d.ts',
          '**/*.config.*',
          '**/mockData/**',
        ],
      },
    },
    
    server: {
      port: 3001,
      host: true,
      // Proxy для локальной разработки (избегаем CORS)
      proxy: {
        '/api': {
          target: env.VITE_API_URL || 'http://localhost:8090',
          changeOrigin: true,
          secure: false,
          configure: (proxy, _options) => {
            proxy.on('proxyReq', (proxyReq, req, _res) => {
              // Сохраняем оригинальный origin для CORS
              if (req.headers.origin) {
                proxyReq.setHeader('Origin', req.headers.origin);
              }
            });
          },
        },
      },
    },
    
    build: {
      outDir: 'dist',
      sourcemap: mode !== 'production',
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['react', 'react-dom', 'react-router-dom'],
            ui: ['@radix-ui/react-dialog', '@radix-ui/react-dropdown-menu', 'framer-motion'],
            query: ['@tanstack/react-query'],
          },
        },
      },
    },
    
    // Environment variables prefix
    envPrefix: 'VITE_',
  };
});

