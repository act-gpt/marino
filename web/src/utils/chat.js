import short from 'short-uuid'
import env from '@/apiService/ENV'
import axios from 'axios'
import { SSE } from './sse.js'
let diff = 0
let secert = {}

const apiPath = '/open/chat/'
const sha256 = '0123456789abcdefghijklmnopqrstuvwxyz'

const ax = axios.create({
  baseURL: env.api,
  withCredentials: true
})

const setSign = (d) => {
  secert = d
  diff = Date.now() - d.timestamp
}

const build = async (url, method, params, data) => {
  const time = timestamp(url)
  params.nonce = secert.nonce
  try {
    params.sign = await encode({ ...params, timestamp: time, path: url }, secert.sign)
  } catch (err) {}
  const headers = {
    'Content-Type': 'application/json',
    'X-Timestamp': time,
    Authorization: secert.sign
  }
  return {
    url: url,
    method: method,
    params,
    data,
    headers
  }
}
const fetch = (url, method, params = {}, data = {}) => {
  return new Promise(async (resolve, reject) => {
    const item = await build(url, method, params, data)
    ax(item)
      .then((res) => {
        if (res.status === 200 && res.data) {
          if (res.data.statusCode === 509) {
            return reject(res)
          }
          resolve(res.data)
        } else {
          reject(res)
        }
      })
      .catch((e) => {
        reject(e)
      })
  })
}

const hash = async (message, secret, algorithm = { name: 'HMAC', hash: 'SHA-256' }) => {
  // encode as UTF-8
  const enc = new TextEncoder('utf-8')
  const key = await crypto.subtle.importKey('raw', enc.encode(secret), algorithm, false, ['sign', 'verify'])
  const signature = await crypto.subtle.sign(algorithm.name, key, enc.encode(message))
  // convert ArrayBuffer to Array
  const hashArray = Array.from(new Uint8Array(signature))

  const hashHex = hashArray.map((b) => b.toString(16).padStart(2, '0')).join('')
  return hashHex
}

const marshall = (params) => {
  params = params || {}
  var keys = Object.keys(params).sort()
  var obj = {}
  var kvs = []
  for (var i = 0; i < keys.length; i++) {
    var k = keys[i]
    if (params[k] === undefined || params[k] === null) {
      delete params[k]
      continue
    }
    obj[k] = params[k]
    kvs.push(keys[i] + '=' + params[k])
  }
  return kvs.join('&')
}

const encode = async (str, secret) => {
  if (Object.prototype.toString.call(str) === '[object Object]') {
    str = marshall(str)
  }
  return await hash(str, secret)
}

const sign = (id) => {
  return new Promise(async (resolve, reject) => {
    try {
      const res = await fetch(apiPath + 'sign/' + id, 'GET')
      const { success, data } = res
      if (success) {
        secert = data
        diff = Date.now() - data.timestamp
      }
      resolve(res)
    } catch (e) {
      reject(e)
    }
  })
}

const sse = async (id ,data, params = {}) => {
  const url = apiPath + 'query/' + id
  //const url = 'http://0.0.0.0:3030/open/chat/query/' + id
  const item = await build(url, 'POST', params, data)
  item.payload = JSON.stringify(item.data)
  item.headers['Accept'] = 'text/event-stream'
  item.headers['Cache-Control'] = 'no-transform'
  var source = SSE(
    env.api +
      url +
      '?' +
      Object.keys(item.params)
        .map((key) => key + '=' + encodeURIComponent(item.params[key]))
        .join('&'),
    item
  )
  return source
}

const bot = (id) => {
  return fetch(apiPath + 'bot/' + id, 'GET')
}

const query = (id, data) => {
  return fetch(apiPath + 'query/' + id, 'POST', {}, data)
}

const datail = (id, doc) => {
  return fetch(apiPath + id, 'GET', { id: doc })
}

const messages = (id, params) => {
  return fetch(apiPath + 'conversation/' + id, 'GET', params)
}

const getId = (str) => {
  var id = 0
  for (var i = 0; i < str.length; i++) {
    id += str.charCodeAt(i)
  }
  return id
}

const timestamp = (str, chars = sha256) => {
  // 所有字符
  const c = chars.split('')
  const id = getId(str)
  // 取特定字符
  const s = c[id % c.length]
  // 特定字符里的次数
  const count =
    str
      .toLowerCase()
      .split('')
      .reduce((curr, text) => (text === s ? ++curr : curr), 0) || 0
  // 服务器当前时间
  const date = Date.now() - diff
  // 计算时间戳
  return date - (count ? date % count : count) + count
}

const parseText = (text) => {
  const items = text.split('\n')
  const reg = /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)/g
  return items
    .reverse()
    .filter((v) => {
      return v.trim()
    })
    .map((text, i) => {
      text = text.trim()
      if (!text) {
        return '<br>'
      }
      if (/^@@/.test(text)) {
        text = text.substr(2, text.length).trim()
        return `<p inner-link="ok"><a href='javascript:void' data-type="send" data-href='${text}'>${text}</a></p>`
      }
      const matchs = text.match(reg)
      if (matchs) {
        matchs.forEach((match) => {
          text = text.replace(match, (m, i) => `<a href="/link?target=${i}" target="_blank">${i}</a>`)
        })
      }
      return `<p>${text}</p>`
    })
    .reverse()
    .join('\n')
}

const conversation = (dev) => {
  const name = dev ? 'devid' : 'cid'
  let conversation = localStorage.getItem(name)
  if (!conversation) {
    conversation = short.generate()
    localStorage.setItem(name, (dev ? 'Dev.' : 'CW.') + conversation)
  }
  return conversation
}

const split = (msgs) => {
  const items = []
  msgs.forEach((item) => {
    if ((item.question || '').trim()) {
      items.push({
        content: item.question,
        direction: 1
      })
    }
    if ((item.answer || '').trim()) {
      items.push({
        content: item.answer,
        direction: 0
      })
    }
  })
  return items
}

const params = (str, key) => {
  const search = str.replace(/[^?]*\?/, '')
  const kv = search.split('&')
  const params = {}
  for (var i = 0; i < kv.length; i++) {
    const sp = kv[i].split('=')
    if (sp[1] && sp[1].trim()) {
      params[sp[0]] = decodeURIComponent(sp[1])
    }
  }
  return key ? params[key] : params
}

const exports =  {
  conversation,
  parseText,
  params,
  timestamp,
  split,
  bot,
  query,
  datail,
  messages,
  setSign,
  sign,
  sse
}

 
export default exports