import cn from 'classnames'

import './style.scss'

import { ReactComponent as OpenLauncher } from './assets/launcher_button.svg'
import { ReactComponent as Close } from './assets/clear-button.svg'

function Launcher({ toggle, chatId, openLabel, closeLabel, showChat }) {
  const toggleChat = () => {
    toggle()
  }
  return (
    <button type='button' className={cn('__agw-launcher', { '__agw-hide-sm': showChat })} onClick={toggleChat} aria-controls={chatId}>
      {showChat ? <Close className='__agw-close-launcher' alt={closeLabel} /> : <OpenLauncher className='__agw-open-launcher' alt={openLabel} />}
    </button>
  )
}

export default Launcher
