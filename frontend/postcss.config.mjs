const config = {
  plugins: [
    "@tailwindcss/postcss",
    [
      "postcss-preset-env",
      {
        stage: 2,
        features: {
          "oklab-function": true,
          "color-mix": true,
          "custom-properties": false, // Don't transform CSS variables
        },
        autoprefixer: {
          flexbox: "no-2009",
          grid: "autoplace",
        },
        browsers: [
          "> 0.5%",
          "last 2 versions",
          "Firefox ESR",
          "not dead",
          "not IE 11",
        ],
        // Preserve custom properties for theme variables
        preserve: true,
      },
    ],
    "autoprefixer",
  ],
};

export default config;
