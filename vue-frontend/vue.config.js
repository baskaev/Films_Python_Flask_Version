const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8080,  // Порт, на котором будет работать ваш фронтенд
    proxy: { // это проксирование будет работать только во время разработки, в продакшене оно не сработает, т.к. не будет devServer 
      '/api': {
        target: 'http://backend:8081',  // Порт вашего сервера Go
        changeOrigin: true,
        secure: false,  // Если используете HTTP вместо HTTPS
      },
    },
  },
})
