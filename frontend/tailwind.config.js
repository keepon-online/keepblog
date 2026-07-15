/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        'primary': 'var(--el-color-primary)',
        'bg-color': 'var(--el-bg-color)',
        'text-color-primary': 'var(--el-text-color-primary)',
        'text-color-regular': 'var(--el-text-color-regular)'
      }
    }
  },
  plugins: [],
  darkMode: 'class'
}