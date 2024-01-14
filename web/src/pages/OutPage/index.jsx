import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import Component from '@/components/Chat'
import './style.scss'
import '@/style/App.scss'
import logo from '@/assets/images/logo.png'

const Page = (props) => {
  const { t } = useTranslation()
  const [bot, setBot] = useState({})
  const callback = (bot) => {
    setBot(bot)
    document.title =  t('title') + (' - ' +bot.setting?.name ? bot.setting?.name  : '')
  }
  return (
    <>
      <div className='d-flex flex-column h-100'>
        <div className='py-2 text-center header'>
          <h4>{bot.setting?.name}</h4>
        </div>
        <div className='main max_width'>
          <div className='max_width chat'>
            <Component {...props} callback={callback} />
          </div>
          <div className='text-center footer sub-title max_width'>
            <a target="_blank"  rel="noreferrer" href='//act-gpt.com/?utm_source=chat'> <img src={logo} alt="logo" className="mr-2 grayColor"
                style={{ width: '24px', height: '24px' }} /> Powered by ACT GPT</a>
          </div>
        </div>
      </div>
    </>
  )
}
export default Page
