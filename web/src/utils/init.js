import Api from '@/apiService/Api'
import i18n from 'i18next'

const init = (props) => {
  const req = async () => {
    let res
    try {
      res = await Api.me()
    } catch (e) {
      return { permission: false, message: 'Plaese sign in' }
    }
    if (res.success) {
      const user = res.data
      props.setUser(user)
      const r = await Api.bots()
      const bots = r.data
      if (!bots || !bots.length) {
        props.setApp({
          auth: true,
          language: i18n.language
        })
        return { permission: true, resource: false, lang: i18n.language }
      }
      const m = props.location.pathname.match(/\/admin\/([^/]+)\//i)
      const id = (m && m[1]) || null
      const bot = bots.find((t) => t.isCurrent || t.id === id) || bots[0]
      if (!bot) {
        //return { permission: false, resource: false, lang: i18n.language }
      }
      props.setBot(bot)
      // auth, notfound
      props.setApp({
        auth: true,
        language: i18n.language,
        notfound: false,
        bots
      })
      return { permission: true, lang: i18n.language }
    }
    props.setApp({
      auth: false,
      language: i18n.language
    })
    return { permission: false, message: 'Plaese sign in' }
  }

  return req()
}

export default init