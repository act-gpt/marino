const { override, addWebpackAlias, addDecoratorsLegacy, adjustStyleLoaders } = require('customize-cra')
const path = require('path')

const multipleEntry = require('react-app-rewire-multiple-entry')([
    {
        entry: 'src/embed.js',
        template: 'public/embed.html',
        outPath: '/embed.html'
    },
    {
        entry: 'src/chat.js',
        template: 'public/chat.html',
        outPath: '/chat.html'
    },
]);

const myOverrides = (config) => {
    return config
}
module.exports = override(
    (config) => {
        config.optimization.splitChunks = {
          cacheGroups: { default: false }
        };
        config.optimization.runtimeChunk = false;
        return config;
    },
    multipleEntry.addMultiEntry,
    addWebpackAlias({
        ['@']: path.resolve(__dirname, 'src')
    }),
    addDecoratorsLegacy(),
    adjustStyleLoaders(({ use }) => {
        use.forEach((loader) => {
          if (/mini-css-extract-plugin/.test(loader.loader)) {
            loader.loader = require.resolve('style-loader');
            loader.options = {};
          }
        });
    }),
    myOverrides,
)