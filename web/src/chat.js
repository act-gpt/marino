import React from 'react'
import { createRoot } from 'react-dom/client'
import utils from '@/utils/chat'
import Chat from './pages/OutPage'
import i18n from './i18n'
import '@/style/base.scss'
import '@/style/utils.scss'

const location = window.location
const path = location.pathname
const params = utils.params(location.search)

let id = params.id || path.replace('/chat/', '').replace(/\/$/, '')
i18n.then((t) => {
  const root = createRoot(document.getElementById('root'))
  root.render(<Chat {...{ id, user: params.user, token: params.token, input: params.q, dev: false }} />)
})
