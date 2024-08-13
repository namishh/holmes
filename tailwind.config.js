/** @type {import('tailwindcss').Config} */

const colors = require("tailwindcss/colors");
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extend: {
      ...colors,
    },
    fontFamily: {
      sans: ["Inter"],
    },
  },
  plugins: [],
};
