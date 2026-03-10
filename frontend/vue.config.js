const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: './',
  configureWebpack: {
    resolve: {
      fallback: {
        crypto: false,
        stream: false,
        assert: false,
        http: false,
        https: false,
        os: false,
        url: false,
        zlib: false
      }
    }
  },
  lintOnSave: false,
  chainWebpack: config => {
    config.plugin('define').tap(args => {
      args[0].__VUE_PROD_DEVTOOLS__ = false
      return args
    })
  }
})
