import React, { useState } from 'react'
import { createRoot } from 'react-dom/client'
import i18nnext from 'i18next'
import i18n from './i18n'
import App from './App'

import 'moment/locale/zh-cn'
import 'moment/locale/zh-tw'
import 'moment/locale/ja'
import { Provider } from 'react-redux'
import store from '@/redux/store'

import reportWebVitals from './reportWebVitals'

i18n.then((t) => {
  const root = createRoot(document.getElementById('root'))
  root.render(
    <Provider store={store}>
      <App lang={i18nnext.language} />
    </Provider>
  )
})

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals(console.log)
