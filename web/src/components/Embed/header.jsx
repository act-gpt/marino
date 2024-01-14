import { ReactComponent as Close } from './assets/clear-button.svg'

function Header({ title, subtitle, toggleChat, showCloseButton, titleAvatar }) {
  return (
    <div className='__agw-header'>
      {showCloseButton && (
        <button className='__agw-close-button' onClick={toggleChat}>
          <Close className='__agw-close' alt='close' />
        </button>
      )}
      <h4 className='__agw-title'>
        {titleAvatar && <img src={titleAvatar} className='avatar' alt='profile' />}
        {title}
      </h4>
      <span>{subtitle}</span>
    </div>
  )
}

export default Header
