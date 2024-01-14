import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import actions from '@/redux/actions'
import init from '@/utils/init'
import { message, Layout, Row, Col, Card, Image, Menu, Modal, Alert, Button, Form, Result, Dropdown, Avatar, FloatButton} from 'antd'
import { useLocation, useHistory } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { PlusOutlined, CheckOutlined, DeleteOutlined, ExclamationCircleOutlined, DownloadOutlined, CustomerServiceOutlined, QuestionCircleOutlined, CommentOutlined} from '@ant-design/icons'
import avatar from '@/assets/images/2021716-105040.jpeg'
import CommonModal from '@/commonComponents/commonModal'
import Contact from '@/components/Contact'
import Loading from '@/components/Loading'
import Api from '@/apiService/Api'
import './style.scss'
import { createFromIconfontCN } from '@ant-design/icons'

import logo from '@/assets/images/logo.png'

const { Header, Content } = Layout
const IconFont = createFromIconfontCN({
  scriptUrl: '/js/icon.js'
})

const Index = (props) => {
  const { user } = props
  const { t, i18n } = useTranslation()
  const location = useLocation()
  const history = useHistory()
  const [apps, setApps] = useState([])
  const [bots, setBots] = useState([])
  const [org, setOrg] = useState({})
  const [app, setApp] = useState(null)

  const [loading, setLoading] = useState(true)
  const [appLoaded, setAppLoaded] = useState(false)
  const [visibleCreate, setVisibleCreate] = useState(false)
  const [visibleContact, setVisibleContact] = useState(false)
  const [visibleInstall, setVisibleInstall] = useState(window.deferredPrompt ? true : false)
  const [banner, setBanner] = useState("")
  const [created, setCreated] = useState(null)
  const [modal, contextHolder] = Modal.useModal()
  const [form] = Form.useForm()

  const installprompt = (event) => {
    const host = process.env.REACT_APP_WEB || 'act-gpt.com|online-gpt.net'
    if(localStorage.getItem('install.alert') || !(new RegExp(host)).test(window.location.href)){
      return
    }
    if(event){
      event.preventDefault()
      console.log('👍', 'beforeinstallprompt', event)
      window.deferredPrompt = event
    }
    window.deferredPrompt && setVisibleInstall(true)
  }
  window.addEventListener('beforeinstallprompt', installprompt)

  window.addEventListener('appinstalled', (event) => {
    console.log('👍', 'appinstalled', event);
    setVisibleInstall(false)
    window.deferredPrompt = null;
  });

  const install = async () => {
      console.log('👍', 'butInstall-clicked')
      const promptEvent = window.deferredPrompt;
      if (!promptEvent) {
        // The deferred prompt isn't available.
        return;
      }
      // Show the install prompt.
      promptEvent.prompt();
      console.log("show")
      // Log the result
      const result = await promptEvent.userChoice;
      console.log('👍', 'userChoice', result);
      // Reset the deferred prompt variable, since
      // prompt() can only be called once.
      window.deferredPrompt = null;
  }

  const getApps = async () => {
    const res = await Api.apps(i18n.language)
    const { success, data } = res
    if (success) {
      setApps(data)
    }
    setAppLoaded(true)
  }

  const getOrg = async () => {
    const res = await Api.org()
    const { success, data } = res
    if (success) {
      setOrg(data)
    }
  }

  const create = (app) => {
    if (!app) {
      message.error(t('home.template_required'))
      return
    }
    const data = {
      id: app.id
    }
    Api.add_bot(data)
      .then(async (res) => {
        const { success, data, code } = res
        if (!success) {
          return message.error(t(`codes.${code}`))
        }
        const items = [...bots]
        items.unshift(data)
        props.setApp({
          auth: true,
          language: i18n.language,
          notfound: false,
          bots: items
        })
        setCreated(data)
        setBots(items)
        //setVisibleCreate(false)
      })
      .catch((err) => {
        console.error(err)
      })
  }

  const deleteBot = async (id) => {
    const confirmed = await modal.confirm({
      title: t('home.confirm.title'),
      icon: <ExclamationCircleOutlined />,
      content: t('home.confirm.message'),
      okText: t('home.confirm.ok'),
      cancelText: t('home.confirm.cancel'),
      onOk: async () => {
        const res = await Api.delete_bot(id)
        const { success, code, data } = res
        if (!success) {
          return message.error(t('codes.' + code))
        }
        const items = bots.filter((app) => app.id != data.id)
        props.setApp({
          bots: items
        })
        setBots(items)
      }
    })
  }

  const onFinish = (values) => {
    if (!app) {
      return message.error(t('home.template_required'))
    }
    create(app)
  }

  const onFinishFailed = (err) => {
    console.error('err', err)
  }

  const changeLanguage = (lng) => {
    history.replace(location.pathname + '?lang=' + lng)
    window.location.reload()
  }

  const loginOut = async () => {
    localStorage.removeItem('user')
    await Api.signout()
    history.push('/')
  }

  
  const languages = [
    { key: 'en', name: 'English' },
    { key: 'zh-CN', name: '简体中文' }
  ]
  
  useEffect(() => {
    Api.banner(props.app.language).then((res)=>{
      setBanner(res)
    }).catch(()=>{})
  },[])
 
  const menu = (
    <Menu style={{ borderRadius: 10, overflow: 'hidden' }}>
      <Menu.Item key={Math.random()} className='flex-column align-items-start' style={{ width: '100%', background: '#3370ff', display: 'flex', padding: 10 }}>
        {user.id ? (
          <>
            <span style={{ display: 'block', color: '#fff' }}>{user.display_name}</span>
          </>
        ) : (
          <span style={{ display: 'block', color: '#fff' }}>Loading...</span>
        )}
      </Menu.Item>
      <Menu.SubMenu title={t('header.language')}>
        {languages.map((lang) => (
          <Menu.Item key={lang.key} style={{ paddingLeft: 10, paddingRight: 10, width: 120 }} onClick={() => changeLanguage(lang.key)}>
            {' '}
            {lang.name}{' '}
            {lang.key == props.app.language ? (
              <span className='pl-4 mainColor'>
                <CheckOutlined />{' '}
              </span>
            ) : (
              ''
            )}
          </Menu.Item>
        ))}
      </Menu.SubMenu>
      <Menu.Divider />
      <Menu.Item>
        <span style={{ display: 'inline-block', width: '100%', color: 'red' }} onClick={loginOut}>
          {t('header.logout')}
        </span>
      </Menu.Item>
    </Menu>
  )

  useEffect(() => {
    // 单点登录
    init(props).then(async (res) => {
      if (!res.permission) {
        await Api.signout()
        return history.push('login?redirect=' + decodeURIComponent(location.pathname))
      }
      setLoading(false)
      getOrg()
      getApps()
    }).catch(()=> setLoading(false))
    installprompt()
  }, [])

  useEffect(() => {
    if (props.app.bots?.length) setBots(props.app.bots)
  }, [props.app.bots])

  return (
    <>
      <Layout className='app' style={{ minHeight: '100vh' }}>
        <Header className='header' style={{ display: 'block' }}>
          <Row>
            <Col flex='1 1 200px'>
            <div className='d-flex align-items-center'>
            <Image
              className="mr-3"
              style={{ width: '28px', height: '28px' }}
              src={logo}
              preview={false}/>
              <h3>{org.name}</h3>
            </div>
            </Col>
            <Col flex='0 1 60px'>
              <Dropdown placement='bottomRight' overlay={menu} overlayStyle={{ width: '20rem' }}>
                <div className='ant-dropdown-link'>
                  <Avatar src={user.avatar || avatar} alt='avatar' style={{ cursor: 'pointer' }} />
                </div>
              </Dropdown>
            </Col>
          </Row>
        </Header>
        {
          visibleInstall && (
            <Alert
              className="install-alert"
              message={t('home.install')}
              description={t('home.install_desc')}
              type="info"
              showIcon
              afterClose={() => {
                localStorage.setItem('install.alert', 1)
              }}
              action={
                <div className="d-flex justify-content-center">
                  <Button icon={<DownloadOutlined />} shape="round" className="mr-3" size="small" type="primary" onClick={install}>
                  {t('home.install')}
                  </Button>
                </div>
              }
              closable
            />
          )
        }
        {banner && !localStorage.getItem(banner.id) ?  (banner.img ? <a href={banner.href} target="_blank"><Image src={banner.img} width={banner.width || '100%'} /> </a> : <Alert showIcon message={banner.message} type={banner.type || "info"} closable onClose={()=>{localStorage.setItem(banner.id, 1)}}/>) : "" }
        <Content style={{ backgroundColor: '#eee', paddingBottom: '60px' }}>
          {
            loading || !appLoaded ? <div  style={{marginTop:"120px"}}><Loading /></div> : (
              <div className='app-manage'>
              <div className='my-apps px-5 pt-5'>
                <h3 className='py-3'>{t('home.myapp')}</h3>
                <Row gutter={[16, 16]}>
                  {contextHolder}
                  {bots.map((bot) => {
                    const items = [
                      {
                        label: <a href={`/admin/${bot.id}/knowledges`}>{t('menu.faq')}</a>,
                        key: 'knowledges',
                        icon: <IconFont type={'icon-knowledge'} />
                      },
                      {
                        label: <a href={`/admin/${bot.id}/chat`}>{t('menu.chat')}</a>,
                        key: 'chat',
                        icon: <IconFont type={'icon-chat'} />
                      },
                      {
                        label: <a href={`/admin/${bot.id}/setting`}>{t('menu.setting')}</a>,
                        key: 'setting',
                        icon: <IconFont type={'icon-setting'} />
                      },
                      {
                        label: <a href={`/admin/${bot.id}/messages`}>{t('menu.messages')}</a>,
                        key: 'messages',
                        icon: <IconFont type={'icon-history'} />
                      },
                      {
                        label: (
                          <a
                            onClick={(e) => {
                              deleteBot(bot.id)
                            }}
                            href='javascript:void(0)'
                          >
                            {t('home.delete')}
                          </a>
                        ),
                        key: 'delete',
                        icon: <DeleteOutlined />
                      }
                    ]
                    return (
                      <Col span={8} key={bot.id}>
                        <Card className='card'>
                          <h6 className='pb-3'>{bot.setting?.name}</h6>
                          <div className='pb-3 sub-title text-truncate'>{bot.setting?.description}</div>
                          <Menu mode='horizontal' style={{ borderBottom: 'none', marginLeft: "-1.3rem"}}>
                            {items.map((item) => {
                              return (
                                <Menu.Item key={item.key}>
                                  {item.icon}
                                  <span style={{ textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>{item.label}</span>
                                </Menu.Item>
                              )
                            })}
                          </Menu>
                        </Card>
                      </Col>
                    )
                  })}
                  <Col span={8} key={'add'} style={{ minHeight: '168px' }}>
                    <Card className='card'>
                      <div className='h-100 text-center d-flex justify-content-center align-items-center'>
                        <Button
                          type='text'
                          icon={<PlusOutlined />}
                          style={{ marginRight: 10 }}
                          onClick={() => {
                            setVisibleCreate(true)
                          }}
                        >
                          {t('home.create_desk')}
                        </Button>
                      </div>
                    </Card>
                  </Col>
                </Row>
              </div>
            </div>
            )
          }
          
        </Content>
        <CommonModal
          visible={visibleCreate}
          title={
            <h3 className='mt-2' style={{ fontWeight: 600 }}>
              {t('home.welcome')}
            </h3>
          }
          common_cancel={() => {
            const items = apps.map((app) => {
              app.selected = false
              return app
            })
            setApps([...items])
            setApp(null)
            setCreated(null)
            setVisibleCreate(false)
          }}
          width={800}
          footer={false}
          custom_foot={
            created ? (
              <></>
            ) : (
              <div className='d-flex justify-content-between mt-5 mb-3'>
                <Button
                  className='px-5'
                  type='primary'
                  onClick={() => {
                    form.submit()
                  }}
                >
                  {t('home.create')}
                </Button>
              </div>
            )
          }
        >
          {created ? (
            <>
              <Result
                status='success'
                title={t('home.result_title')}
                subTitle={t('home.result_description')}
                extra={[
                  created.setting.model ? (
                    <Button
                      key='chat'
                      onClick={() => {
                        props.history.push(`/admin/${created.id}/chat`)
                      }}
                    >
                      {t('menu.chat')}
                    </Button>
                  ) : (
                    <Button
                      key='knowledges'
                      onClick={() => {
                        props.history.push(`/admin/${created.id}/knowledges`)
                      }}
                    >
                      {t('menu.faq')}
                    </Button>
                  ),
                  <Button
                  type='primary'
                  key='setting'
                  onClick={() => {
                    props.history.push(`/admin/${created.id}/setting`)
                  }}>
                  {t('menu.setting')}
                </Button>
                ]}
              />
            </>
          ) : (
            <Form form={form} layout='vertical' requiredMark={false} onFinish={onFinish} onFinishFailed={onFinishFailed}>
              {/*
              <Form.Item name='name' label={<strong>{t('home.desk_name')}</strong>} rules={[{ required: true, message: '' }]}>
                <Input autoFocus />
              </Form.Item>
              */}
              <h5 className='pb-2'>{t('home.desk_template')}</h5>
              <Form.Item noStyle>
                <Row gutter={[8, 8]} className='app-templates'>
                  {apps.map((app) => {
                    return (
                      <Col span={8} key={app.id}>
                        <a
                          href='#'
                          onClick={(e) => {
                            e.preventDefault()
                            const items = apps.map((app) => {
                              app.selected = false
                              return app
                            })
                            app.selected = true
                            setApp(app)
                            setApps([...items])
                          }}
                        >
                          <Card style={{ borderRadius: '6px' }} className={app.selected ? 'selected' : ''}>
                            <h6 className='pb-3'>{app.name}</h6>
                            <div className='pb-3 sub-title text-truncate'>{app.description}</div>
                          </Card>
                        </a>
                      </Col>
                    )
                  })}
                </Row>
              </Form.Item>
            </Form>
          )}
        </CommonModal>
      </Layout>
      <Contact
        show={visibleContact}
        onChange={(show) => {
          setVisibleContact(show)
        }}
        title={t('contact_title')}
        description={t('contact_description')}
      />
      <FloatButton.Group
        trigger="hover"
        type="primary"
        style={{ right: 24 }}
        icon={<CustomerServiceOutlined />}>
        <FloatButton  tooltip={<div>{t('home.open_document')}</div>} href={process.env.REACT_DOC_PATH || "/doc"} target="_blank"/>
        {props.app.language == "zh-CN" 
        ? <FloatButton  tooltip={<div>{t('home.contact_us')}</div>} icon={<CommentOutlined />} onClick={()=> setVisibleContact(true)} />
        : <FloatButton  tooltip={<div>{t('discord')}</div>} href={"https://discord.gg/GJrmRDyh"} target="_blank"/>}
      </FloatButton.Group>
    </>
  )
}

const mapStateToProps = (state) => {
  return {
    user: state.user,
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}
export default connect(mapStateToProps, mapDispatchToProps)(Index)
