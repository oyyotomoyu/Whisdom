import path from "node:path"
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  define: {
    "process.env.NODE_ENV": JSON.stringify(process.env.NODE_ENV ?? "development"),
  },
  resolve: {
    alias: {
      // Branding (colors, favicon, login logo) lives outside client/ so an
      // ODM build can replace it without touching component code. See
      // docs/client.md's "Theme And Branding (ODM)" section.
      "@odm": path.resolve(import.meta.dirname, "../odm"),
    },
  },
  server: {
    fs: {
      allow: [path.resolve(import.meta.dirname, "..")],
    },
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
})
