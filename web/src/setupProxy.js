const { createProxyMiddleware } = require('http-proxy-middleware')
const fs = require('fs')
const axios = require("axios")
module.exports = function (app) {
  const APP_HOST = process.env.APP_HOST
  const url = `http://${APP_HOST ? APP_HOST : '127.0.0.1'}:6789/`
  app.use((req, res, next) =>{
      res.set('X-Server', 'web')
      next()
  });

  app.use(`/chat/:id`, (req, res) =>{
    res.set('Content-Type', 'text/html')
    const url = "http://0.0.0.0:" + (process.env.PORT || 3000) + "/chat.html"
    axios({
      method: 'get',
      url,
      responseType: 'stream'
    }).then(function (response) {
      response.data.pipe(res)
    }).catch((e) => {
      res.send(Buffer.from("Server error" + e))
    })
  })
  app.use(
    '/api',
    createProxyMiddleware({
      target: url,
      changeOrigin: true
    })
  )
  /*
  app.use(
    '/open/chat/query',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
      selfHandleResponse: true, 
      on: {
        onProxyRes(responseBuffer, proxyRes, req, res) {
          if (req.headers.accept === 'text/event-stream') {
            res.writeHead(res.statusCode, res.headers);
            proxyRes.pipe(res);
          }
        }
      }
    })
  )
  */
  app.use(
    '/open',
    createProxyMiddleware({
      target: url,
      changeOrigin: true,
    })
  )
  app.use(
    '/dashboard',
    createProxyMiddleware({
      target: url,
      changeOrigin: true
    })
  )
  app.use(
    '/v1',
    createProxyMiddleware({
      target: url,
      changeOrigin: true
    })
  )
}
