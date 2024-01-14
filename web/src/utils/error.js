import { message } from 'antd'

// eslint-disable-next-line import/no-anonymous-default-export
export default function (props, errResponse) {
  if (errResponse.statusCode === 509) {
    message.error(JSON.stringify(errResponse.message), 1)
    const temp = Object.keys(props)
    if (temp.length !== 0) {
      setTimeout(() => {
        props.history.push('/login')
        localStorage.clear() // 清除缓存
      }, 1100)
      return
    } else {
      return
    }
  } else {
    message.error(errResponse.message)
  }
}
