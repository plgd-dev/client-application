const path = require('path')
// const { ProvidePlugin } = require('webpack')
/* eslint-disable */
const { CracoAliasPlugin } = require('react-app-alias-ex')

module.exports = {
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src/'),
      shared: path.resolve(__dirname, '../shared/src/'),
      sharedOut: path.resolve(__dirname, '../../shared/src/'),
    },
    // plugins: [
    //   new ProvidePlugin({
    //     React: 'react',
    //   }),
    // ],
  },
  plugins: [
    {
      plugin: CracoAliasPlugin,
      options: {},
    },
  ],
}
