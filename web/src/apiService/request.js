import axios from 'axios'
import { message } from 'antd'

axios.interceptors.response.use(
  function (res) {
    // jwt token invalid
    if (res.data && res.data.statusCode === 509) {
      message.error('User authentication failed, please reload page.')
    }
    return res
  },
  function (error) {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    // Do something with response error
    return error
  }
)

function id () {
  function uuidv4() {
    return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c =>
      (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
  }
  let id = localStorage.getItem("id")
  if (id) {
    return id
  }
  id = uuidv4()
  localStorage.setItem("id", id)
  return id
}
export function request(url, method, params = {}, data = {}) {
  return new Promise((resolve, reject) => {
    const user = JSON.parse(localStorage.getItem('user'))
    let access_token = params.access_token
    if (!access_token) {
      delete params.access_token
    } else {
      access_token = user ? user.access_token : ''
    }

    axios({
      url: url,
      method: method,
      params,
      data,
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization'
        Authorization: access_token,
        UUID: id () 
      }
    })
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

export function fetch(url, method, headers = {}, params = {}, data = {}) {
  return new Promise((resolve, reject) => {
    const user = JSON.parse(localStorage.getItem('user'))
    let access_token = params.access_token
    if (!access_token) {
      delete params.access_token
    } else {
      access_token = user ? user.access_token : ''
    }
    headers = {
      'Content-Type': 'application/json',
      // 'Authorization'
      Authorization: access_token,
      ...headers
    }
    axios({
      url: url,
      method: method,
      params,
      data,
      headers
    })
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
