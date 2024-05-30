/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./api/public/**/*.{templ,js}"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["night"],
  },
  plugins: [
    require('daisyui'),
  ],
}

