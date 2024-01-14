import React, { useState, useEffect, useRef} from 'react'
import { Avatar, Watermark, notification } from 'antd'
import { useTranslation } from 'react-i18next'
import { ExclamationOutlined } from '@ant-design/icons'
import ChatUtil from '@/utils/chat'
import CommonModal from '@/commonComponents/commonModal'
import Loading from '@/components/Loading'
import { MainContainer, ChatContainer, MessageList, Message, MessageInput } from '@chatscope/chat-ui-kit-react'
import ChatMessage from './Message'

import '@chatscope/chat-ui-kit-styles/dist/default/styles.min.css'
import '@/style/chat.scss'

const Component = (props) => {
  const { id, user, dev, input, callback } = props
  // conversation id
  const conversation = ChatUtil.conversation(false)
  const { t } = useTranslation()
  // history params
  const [params, setParams] = useState({
    page: 1,
    size: 3
  })
  const [loading, setLoading] = useState(true)
  const [query, setQuery] = useState(input)
  const [bot, setBot] = useState({})
  const [detail, setDetail] = useState('')
  const [show, setShow] = useState(false)
  const [disabled, setDisabled] = useState(false)
  const [key, setKey] = useState(Math.random())
  // messages
  const [messages, setMessages] = useState([])
  const [message, setMessage] = useState(null)
  const msgListRef = useRef()
  const [api, contextHolder] = notification.useNotification();

  const onError = (data,  message) => {
    var content = ""
    if (data.code > 7000 && data.code < 7005){
      content = t('codes.' + data.code)
      message.content = content
    }else{
      content = t('codes.' + data.code)
      if (!content){
        content = t('codes.500')
      }else{
        content = data.message
      }
      api.warning({
        message: t('chat.reminder'),
        description: content,
      });
    }
  }
  const sse = async (input) => {
    let text = ''

    // message
    const  message = {
      direction: 0,
      role: "",
      streaming: true,
      content: '',
      loading: true,
      id: null,
      source: []
    }

    const msgs = [...messages]
    msgs.push({
      role: "user",
      content: input,
      direction: 1
    })
    setMessages(msgs)
    

    // request done
    const onDone = () => {
      message.content = text
      if (text){
        setMessages([...msgs,  message])
      }else{
        setMessages([...msgs])
      }
      setMessage(null)
    }
    const sse = await ChatUtil.sse(bot.id, {
      conversation,
      user,
      prompt: input,
      stream: true,
      //max_tokens: 1000,
      temperature: 0.1
    }, {
      source: dev ? "dev" : "web" 
    })

    // client side event
    sse.addEventListener("open", (e) => {
      message.loading = false
      setMessage(message)
    })

    sse.addEventListener("readystatechange", (e) => {
      if (e.readyState == sse.CONNECTING) {
        setMessage(message)
      }
      // close http request
      if (e.readyState == sse.CLOSED) {
        setDisabled(false)
        message.done = true
        onDone()
      }
    })
    sse.addEventListener("error", (err) => {
       //http status is not 200 or network error
       console.log(err.data)
      try {
        const data = JSON.parse(err.data)
        onError(data)
      }catch(e){
        console.log(e.data)
      }
      onDone()
    })

    sse.addEventListener("abort",  (e) => {
      onDone()
    })

    // server side evvent
    sse.addEventListener("message", (e) => {
      var data
      try{
        data = JSON.parse(e.data)
      }catch(e){
        message.loading = false
        message.streaming = false
        return setMessage({ ... message })
      }
      const { id, code,  object, created, model, choices, source } = data
      // some thing error
      if (code) {
        onError({code, message: data.message},  message)
        message.loading = false
        message.streaming = false
        return setMessage({ ... message })
      }
      if (! message.id) {
        message.id = id
      }
      if (source){
        message.streaming = false
        message.source = source
      }
      if (choices && choices.length){
        const item = choices[0]
        const {delta, finish_reason, index} = item
        const {content, role } = delta
        if (role && ! message.role){
          message.role = role
          //console.info("answer by", role)
        }
        if (finish_reason){
          if  (finish_reason != 'stop'){
            console.warning("finish by", finish_reason)
          }
          message.streaming = false
        }
        if (content){
          text += content
          message.content += content
          return setMessage({ ... message })
        }
      }
      setMessage({ ... message })
    })

    sse.stream()

  }

  const getBot = async () => {
    await ChatUtil.sign(id)
    try {
      const res = await ChatUtil.bot(id)
      const { success, data } = res
      setLoading(false)
      if (!success) {
        callback && callback({ failed: true })
        return setBot({ failed: true })
      } else {
        setBot(data)
        callback && callback(data)
      }
    } catch (e) {
      callback && callback({ failed: true })
      setLoading(false)
      return setBot({ failed: true })
    }
  }

  const getHistories = async () => {
    if (bot.failed) {
      return
    }
    const { page, size } = params

    if (page < 0) {
      return
    }
    const res = await ChatUtil.messages(bot.id, { page, size, cvs: conversation })
    const { success, data } = res
    if (!success) {
      return
    }
    if (data.length === 0) {
      if (page === 1) {
        setMessages([
          {
            content: bot.setting.welcome,
            direction: 0
          }
        ])
      }
      setParams({
        page: -1,
        size
      })
      return
    }
    const vals = ChatUtil.split(Array.prototype.reverse.apply(data))
    Array.prototype.push.apply(vals, messages)
    vals.push({
      content: bot.setting.welcome,
      direction: 0
    })
    setMessages(vals)
    setParams({
      page: data.length === page ? page + size : -1,
      size
    })
  }

  const scroll = () => {
    msgListRef.current?.scrollToBottom('auto')
  }

  const onClick = (e) => {
    let ele = e.target
    e.preventDefault()
    e.stopPropagation()
    // mayby embed image
    if (ele.nodeName !== 'A' && ele.className.indexOf('embed')) {
      ele = ele.parentNode
    }
    if (ele.nodeName !== 'A') {
      return
    }
    const { type } = ele.dataset
    // type: query, embed, opener
    switch (type) {
      case 'query':
        return onSend(ele.title)
      case 'embed':
        return open(ele.href)
      case 'opener':
        return open(ele.href)
    }
  }

  const is = (ele, className) => {
    return ele.className.indexOf(className) > -1
  }

  const open = (url) => {
    window.open(process.env.REACT_APP_WEB + '/link?target=' + encodeURIComponent(url), '_blank').focus()
  }

  const onSend = async (input) => {
    if (!input.trim()) {
      return
    }
    setDisabled(true)
    input = input.replace(/<[^>]+>/g, '')
    return sse(input)
  }

  useEffect(() => {
    setTimeout(scroll, 500)
    if (messages.length && query) {
      onSend(input)
      setQuery('')
    }
  }, [messages])

  useEffect(() => {
    if (!loading) {
      return
    }
    getBot()
  }, [loading])

  useEffect(() => {
    if (bot.setting) {
      getHistories()
    }
  }, [bot])

  return (
    <div key={key} className={dev ? 'px-5' : ''} style={{height: `calc(100vh - 50px${dev ? ' - 24px' : ''})`}}>
      {contextHolder}
      {loading ? (
        <Loading />
      ) : bot.failed ? (
        <div className='text-center d-flex justify-content-center align-items-center h-100'>
          <div className="pr-3">
            <Avatar size={64} style={{ backgroundColor: '#0677FF' }} icon={<ExclamationOutlined />} />
          </div>
          <div className=''>
            <h2 className='mb-1'>{t('chat.not_found')}</h2>
            <div className='sub-title'>{t('chat.contact')}</div>
          </div>
        </div>
      ) : (
        <Watermark className="h-100" content={t('chat.watermark')} font={{color:'rgba(0, 0, 0, 0.02)',  fontSize:"16"}}>
          <MainContainer className='w-100' style={{paddingBottom:"8px" , paddingLeft: dev ? "24px" : "", paddingRight: dev ? "24px": ""}}>
            <ChatContainer>
              <MessageList ref={msgListRef} typingIndicator={null} loading={false} scrollBehavior='smooth'>
                {messages
                  .filter((msg) => msg.content || msg.html || msg.loading)
                  .map((message, i) => {
                    return (
                      <Message
                        key={i}
                        model={{
                          direction: message.direction,
                          type: 'custom'
                        }}>
                        <Message.CustomContent>
                          <ChatMessage {...message} onClick={onClick} />
                        </Message.CustomContent>
                      </Message>
                    )
                  })}
                {message ? (
                  <Message
                    key={Math.random()}
                    model={{
                      direction: message.direction,
                      type: 'custom'
                    }}
                  >
                    <Message.CustomContent>
                      <ChatMessage {...message} dev={dev} onClick={onClick} />
                    </Message.CustomContent>
                  </Message>
                ) : (
                  ''
                )}
              </MessageList>
              <MessageInput attachButton={false} autoFocus fancyScroll={false} disabled={disabled} onSend={onSend} placeholder={t('chat.input_placeholder')} />
            </ChatContainer>
          </MainContainer>
        </Watermark>
      )}
      <CommonModal
        title={t('chat.detail')}
        width={600}
        footer={false}
        visible={show}
        common_cancel={() => {
          setShow(false)
        }}
        children={
          <div>
            <pre>{detail}</pre>
          </div>
        }
      />
    </div>
  )
}

export default Component
