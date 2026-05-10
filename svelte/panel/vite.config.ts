/// <reference types="vitest" />
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    svelte(),
    tailwindcss()
  ],
  resolve: {
    alias: {
      $lib: resolve('./src/lib')
    }
  },
  test: {
    environment: 'node',
    include: ['src/**/*.test.ts'],
  }
})
