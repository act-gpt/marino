import React, { useState } from 'react'
import Component from '@/components/Chat'
import Launcher from './launcher'
import Header from './header'
/*
.rcw-launcher, .rcw-conversation-container .rcw-header {
  background-color: #356be6 !important;
}
*/
function Embed(props) {
  const { user, id, token, show, color } = props
  const [isShow, setIsShow] = useState(show)
  const [setting, setSetting] = useState({})

  const callback = (bot) => {
    const { setting, failed } = bot
    if (failed) {
      setSetting({
        name: 'Unknow bot',
        description: ''
      })
    }
    if (setting) {
      if (setting.color) {
        document.documentElement.style.cssText = `--main-agw-bg-color: ${setting.color}`
      }
      setSetting(setting)
    }
  }
  const toggle = () => {
    setIsShow(!isShow)
  }
  return (
    <>
      {isShow ? (
        <div className='__agw-conversation-container'>
          <Header title={setting.name || ''} subtitle={setting.description || ''} toggleChat={toggle} />
          <div className='__agw-messages-container'>
            <Component {...props} callback={callback} />
          </div>
        </div>
      ) : (
        ''
      )}
      <Launcher toggle={toggle} chatId={''} openLabel={'Open chat'} closeLabel={'Close chat'} showChat={isShow} />
    </>
  )
}

export default Embed
