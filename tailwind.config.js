/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [".views/**/*.{templ,js}"],
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

