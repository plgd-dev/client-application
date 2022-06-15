const path = require('path')
// const { ProvidePlugin } = require('webpack')

module.exports = {
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src/'),
    },
    // plugins: [
    //   new ProvidePlugin({
    //     React: 'react',
    //   }),
    // ],
  },
}
