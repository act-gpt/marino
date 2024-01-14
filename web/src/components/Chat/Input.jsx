import { KeyboardEvent, MutableRefObject, useCallback, useContext, useEffect, useRef, useState } from 'react'
const Component = (props) => {
  const [content, setContent] = useState('')
  const [isTyping, setIsTyping] = useState < boolean > false
  const handleSend = () => {
    if (messageIsStreaming) {
      return
    }

    if (!content) {
      alert(t('Please enter a message'))
      return
    }

    onSend({ role: 'user', content }, plugin)
    setContent('')
    setPlugin(null)

    if (window.innerWidth < 640 && textareaRef && textareaRef.current) {
      textareaRef.current.blur()
    }
  }
}

export default Component
