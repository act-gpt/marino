import React, { useEffect, useState } from 'react'
import PropTypes from 'prop-types'
import { message, Layout, Menu, Dropdown, Tag, Skeleton, Row, Col, Image, Form, Divider, Popover, Tooltip } from 'antd'
import { HomeOutlined, PlusOutlined } from '@ant-design/icons'
import CustomMenu from '@/components/CustomMenu'
import CommonModal from '@/commonComponents/commonModal'
import { useTranslation } from 'react-i18next'
import queryString from 'query-string'
import { NavLink } from 'react-router-dom'
import { useLocation, useHistory } from 'react-router-dom'
import { connect } from 'react-redux'
import actions from '@/redux/actions'

import Create from '@/components/Create'
import fallback from '@/assets/images/logo.png'

const { Sider } = Layout
const { Item, SubMenu } = Menu

const AppAside = (props) => {
  const { t, i18n } = useTranslation()

  let { menuToggle, menu, navigation, user, bot, app, setBot, setApp } = props

  const parsed = queryString.parse(navigation.location.search)
  const location = useLocation()
  const history = useHistory()

  const [form] = Form.useForm()
  const [visible, setVisible] = useState(false)
  const [newVisible, setNewVisible] = useState(parsed.create)

  const change = (id) => {
    setVisible(false)
    setApp({ loading: true })
    const bot = app.bots.find((t) => t.id == id)
    if (bot) {
      setBot(bot)
      navigation.history.push(`/admin/${bot.id}/agents`)
    }
    setTimeout(() => {
      setApp({ loading: false })
    }, 1000)
  }

  useEffect(() => {}, [app.bots])

  const show = () => {
    setVisible(false)
    setNewVisible(true)
  }

  const hide = () => {
    if (!app.bots || !app.bots.length) {
      return message.error(t('left.must_have_one'))
    }
    setNewVisible(false)
  }

  const lgm = () => (
    <div style={{ borderRadius: 10, overflow: 'hidden', width: '260px' }}>
      <div className='d-flex py-5 px-5' style={{ background: '#326BFF' }}>
        <Image
          style={{ width: '40px', height: '40px' }}
          src={bot.avatar || fallback}
          fallback={fallback}
          preview={false}
        />
        <span className='pl-4' style={{ color: '#fff', lineHeight: '40px' }}>
          {bot.setting?.name}
        </span>
      </div>
      <Divider style={{ margin: '0' }} />
      <div style={{ background: '#fff', cursor: 'pointer', color: '#3370ff' }} className='d-flex p-4 link' onClick={(e) => setVisible(true)}>
        {t('left.switch')}
      </div>
    </div>
  )

  const rules = [{ required: true, message: t('required') }]

  return (
    <React.Fragment>
      <Sider className='aside' collapsed={menuToggle} width={250} style={{ borderRight: '1px solid #f0f0f0' }}>
        <div className='logo ml-2' style={{ display: 'flex' }}>
          <Skeleton loading={app.loading} avatar active paragraph={{ rows: 0 }}>
            <Image
              className="mr-2"
              style={{ borderRadius: '6px', display: 'inline-block', objectFit: 'cover'}}
              src={bot.avatar || fallback}
              alt={bot.name}
              fallback={fallback}
              preview={false}
            />
            <span className='name h4 text-nowrap'> {bot.name}</span>
          </Skeleton>
        </div>
        <CustomMenu menuToggle={menuToggle} menu={menu} />
        <div className='slide_back'>
          <Menu mode='inline'>
            <Menu.Item key={'home'}>
              <NavLink to={'/apps'}>
                <HomeOutlined />
                <span style={{textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>{t('back_to_home')}</span>
              </NavLink>
            </Menu.Item>
          </Menu>
        </div>
        <CommonModal
          title={t('left.select')}
          width={600}
          footer={false}
          visible={visible}
          common_cancel={() => {
            setVisible(false)
          }}
          children={
            <div>
              {visible &&
                app.bots.map((bot) => (
                  <Row key={bot.id} className='hand' style={{ paddingBottom: '10px' }} onClick={(e) => change(bot.id)}>
                    <Col span={4}>
                      <Image src={bot.avatar || fallback} fallback={fallback} preview={false} />
                    </Col>
                    <Col span={20} className='flex'>
                      <div className='flex-column-center align-self-center' style={{ paddingLeft: '20px', alignItems: 'flex-start', height: '100%' }}>
                        <h3>{bot.setting?.name}</h3>
                      </div>
                    </Col>
                  </Row>
                ))}
              <Row className='hand' style={{ paddingBottom: '10px' }} onClick={(e) => show()}>
                <Col span={4}>
                  <i className='' style={{ width: '100%', textAlign: 'center', display: 'inline-block', fontSize: '60px', backgroundColor: '#3370FF', color: '#fff', borderRadius: '10px' }}>
                    <PlusOutlined />
                  </i>
                </Col>
                <Col span={20} className='flex'>
                  <div className='flex-column-center' style={{ paddingLeft: '20px', alignItems: 'flex-start', height: '100%' }}>
                    <h3> {t('left.create')} </h3>
                  </div>
                </Col>
              </Row>
            </div>
          }
        />
        <Create {...props} visible={newVisible} setVisible={setNewVisible} />
      </Sider>
    </React.Fragment>
  )
}

AppAside.propTypes = {
  menuToggle: PropTypes.bool,
  menu: PropTypes.array.isRequired
}

const mapStateToProps = (state) => {
  return {
    user: state.user,
    app: state.app,
    bot: state.bot
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(AppAside)
