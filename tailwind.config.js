/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./api/public/**/*.{templ,js}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}

