var path = require('path');
var webpack = require('webpack');

const fs  = require('fs');
//const lessToJs = require('less-vars-to-js');
//const themeVariables = lessToJs(fs.readFileSync(path.join(__dirname, './static/ant-theme-vars.less'), 'utf8'));

    module.exports = {
      devtool: 'eval',
      entry: [
//        'webpack-dev-server/client?http://localhost:8001',
        path.join(__dirname, 'src', 'index')
      ],
      output: {
        path: path.join(__dirname, 'dist'),
        filename: 'bundle.js',
        publicPath: '/'
      },
      resolve: {
        extensions: ['.js', '.ts', '.tsx', '.json']
      },
      module: {
        loaders: [{
          test: /\.tsx?$/,
          loaders: ['babel-loader?presets[]=react', 'ts-loader'],
          exclude: /node_modules/
        }, {
          test: /\.(jsx?)$/,
          loaders: ['babel'],
          exclude: /node_modules/
        }, {
          test: /\.css$/,
          loaders: [ 'style-loader', 'css-loader', 'sass-loader' ]
        }, {
          test: /\.scss$/,
          loaders: [ 'style-loader', 'css-loader', 'sass-loader' ]
        }, {
          test: /\.less$/,
          use: [{
              loader: "style-loader" // creates style nodes from JS strings
          }, {
              loader: "css-loader" // translates CSS into CommonJS
          }, {
              loader: "less-loader", // compiles Less to CSS
              /*options: {
                modifyVars: themeVariables
              }*/
          }]
        }]
      }
    };