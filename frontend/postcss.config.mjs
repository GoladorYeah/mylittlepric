const config = {
  plugins: [
    "@tailwindcss/postcss",
    [
      "postcss-preset-env",
      {
        stage: 3,
        features: {
          "oklab-function": true,
          "color-mix": true,
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
      },
    ],
  ],
};

export default config;
