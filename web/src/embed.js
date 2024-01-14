import React, { useState } from 'react'
import i18n from './i18n'
import { createRoot } from 'react-dom/client'
import Embed from './components/Embed'
import '@/style/base.scss'
import '@/style/utils.scss'

const params = (str, key) => {
  const search = str.replace(/[^?]*\?/, '')
  const kv = search.split('&')
  const params = {}
  for (var i = 0; i < kv.length; i++) {
    const sp = kv[i].split('=')
    if (sp[1] && sp[1].trim()) {
      params[sp[0]] = sp[1]
    }
  }
  return key ? params[key] : params
}

const scripts = document.scripts
const script = scripts[scripts.length - 1]
const src = script.getAttribute('src')
const items = params(src)
const { id, user } = window.__act_gpt || items
const container = document.createElement('div')
container.id = '__act__gpt__embed'
container.className = '__act__gpt__embed'
document.body.appendChild(container)
const root = createRoot(container)
i18n.then((t) => {
  root.render(<Embed {...{ id: id, user: user }} />)
})
