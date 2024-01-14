import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import loadable from './utils/loadable'
import React, { useState, useEffect } from 'react'

import { ConfigProvider } from 'antd'
import { useTranslation } from 'react-i18next'
import { connect } from 'react-redux'
import moment from 'moment'
import actions from '@/redux/actions'
import CommonModal from '@/commonComponents/commonModal'

import 'animate.css'
import '@/style/base.scss'
import '@/style/App.scss'
import '@/style/layout.scss'
import '@/style/utils.scss'
import '@/style/base_antd.scss'


import enUS from 'antd/lib/locale/en_US'
import zhCN from 'antd/lib/locale/zh_CN'

const locals = {
  en: enUS,
  'zh-Hans': zhCN,
  'zh-CN': zhCN
}

const DefaultLayout = loadable(() => import('./containers/DefaultLayout'))
const Forbidden = loadable(() => import('./pages/403'))
const NotFound = loadable(() => import('./pages/404'))
const Index = loadable(() => import('./pages/Index'))
const Home = loadable(() => import('./pages/Home'))
const Login = loadable(() => import('./pages/Login'))
const Signup = loadable(() => import('./pages/Signup'))
const Welcome = loadable(() => import('./pages/Welcome'))
const Setup = loadable(() => import('./pages/Setup'))

const items = {
  en: 'en',
  'zh-Hans': 'zh-cn'
}

const App = (props) => {
  const { app, setApp } = props
  const language = app.language
  const { t, i18n } = useTranslation()
  const [local, setLocal] = useState(locals[language] || locals['en'])
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const lang = props.lang
    const l = locals[lang]
    // set html element lang attrubite
    document.documentElement.lang = lang
    setLocal(l || locals['en'])
    moment.locale(items[lang])
    i18n.changeLanguage(lang)
    setApp({ ...app, language: lang })
    document.title = t('title')
  }, [])

  return (
    <ConfigProvider 
      theme={{
        token: {
          colorPrimary: '#0677FF',
          colorInfo: '#0677FF',
          colorLink: '#0677FF',
        },
      }}
      locale={local}>
      <Router>
        <Switch>
          <Route key='login' path='/login' exact component={Login} />
          <Route key='signup' path='/signup' exact component={Signup} />
          <Route key='404' path='/404' exact component={NotFound} />
          <Route key='403' path='/403' exact component={Forbidden} />
          <Route key='home' path='/apps' exact component={Home} />
          <Route key='welcome' path='/welcome' exact component={Welcome} />
          <Route key='setup' path='/setup' component={Setup} />
          <Route key='admin' path='/admin/*' component={DefaultLayout} />
          <Route key='index' path='/' exact component={Index} />
          <Route key='not' when={false} exact component={NotFound} />
        </Switch>
      </Router>
      <CommonModal
        title={t('abandon_edit')}
        width={600}
        visible={visible}
        danger={true}
        common_cancel={() => {
          setVisible(false)
        }}
        common_confirm={() => {
          setVisible(false)
          window.navigation.history.push(window.location.pathname)
        }}
        children={<div className='py-4 my-4'>{t('abandon_edit_desc')}</div>}
      />
    </ConfigProvider>
  )
}
const mapStateToProps = (state) => {
  return {
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}
export default connect(mapStateToProps, mapDispatchToProps)(App)
