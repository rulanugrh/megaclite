/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html}"],
  theme: {
    extend: {
      colors: {
        lightColor: '#ECDFCC',
        lightGray: '#697565',
        boldGray: '#3C3D37',
        darkColor: '#1E201E'
      },
      container: {
        center: true,
        padding: "16px"
      },
      fontFamily: {
        jetBrains: "'JetBrains Mono', serif"
      }
    }
  },
  plugins: [],
}

