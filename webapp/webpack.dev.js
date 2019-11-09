const merge = require('webpack-merge');
const common = require('./webpack.common.js');
var path = require('path');

module.exports = merge(common, {
  devtool: 'inline-source-map',
  devServer: {
    port: 9001,
    host: '0.0.0.0',
    contentBase: [ './dist', './' ],
    publicPath: '/',
    historyApiFallback: true
  }
} );
