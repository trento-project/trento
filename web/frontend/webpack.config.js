/* eslint-disable no-undef */
const path = require('path');

const isProduction = process.env.NODE_ENV == 'production';

const stylesHandler = 'style-loader';

const config = {
  entry: {
    check_results: './javascripts/check_results.js',
  },
  output: {
    path: path.resolve(__dirname, 'assets/js'),
    filename: '[name].js',
  },
  plugins: [],
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/i,
        loader: 'babel-loader',
      },
      {
        test: /\.css$/i,
        use: [stylesHandler, 'css-loader'],
      },
      {
        test: /\.s[ac]ss$/i,
        use: [stylesHandler, 'css-loader', 'sass-loader'],
      },
      {
        test: /\.(eot|svg|ttf|woff|woff2|png|jpg|gif)$/i,
        type: 'asset',
      },
    ],
  },
  resolve: {
    alias: {
      '@components': path.resolve(__dirname, 'components'),
    },
  },
};

module.exports = () => {
  if (isProduction) {
    config.mode = 'production';
  } else {
    config.mode = 'development';
  }
  return config;
};
