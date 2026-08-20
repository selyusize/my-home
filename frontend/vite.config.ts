import path from "node:path";
import { fileURLToPath } from "node:url";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

const src = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "src");

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  resolve: {
    alias: {
      "@": src,
      "@app": path.resolve(src, "app"),
      "@pages": path.resolve(src, "pages"),
      "@widgets": path.resolve(src, "widgets"),
      "@features": path.resolve(src, "features"),
      "@entities": path.resolve(src, "entities"),
      "@shared": path.resolve(src, "shared"),
    },
  },
  plugins: [react(), tailwindcss(), wails("./bindings")],
});
