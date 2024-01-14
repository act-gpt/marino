import React, { useState } from 'react'
import PropTypes from 'prop-types'
import { message, Menu, Dropdown, Layout, Avatar, Switch, Tooltip, Skeleton } from 'antd'
import { MenuFoldOutlined, HomeOutlined, MenuUnfoldOutlined, CheckOutlined } from '@ant-design/icons'
import avatar from '@/assets/images/2021716-105040.jpeg'
import { useTranslation } from 'react-i18next'

import { useLocation, useHistory } from 'react-router-dom'

import { connect } from 'react-redux'

import actions from '@/redux/actions'
const { Header } = Layout

const AppHeader = (props) => {
  const { menuClick, menuToggle, loginOut, user, app, setApp } = props
  const { t, i18n } = useTranslation()
  const location = useLocation()
  const history = useHistory()

  const changeLanguage = (lng) => {
    history.replace(location.pathname + '?lang=' + lng)
    window.location.reload()
  }

  const [visible, setVisible] = useState(false)
  const [defaultChecked, setDefaultChecked] = useState(true)
  const change_switch = (checked) => {
    checked ? setVisible(false) : setVisible(true)
  }

  const confirm = () => {
    setVisible(false)
    setDefaultChecked(true)
    message.info('comning soon...')
  }

  const backHome = () => {
    setApp({ loading: true })
    history.push('/apps')
    window.location.reload()
    setTimeout(() => {
      setApp({ loading: false })
    }, 1000)
  }

  const languages = [
    { key: 'en', name: 'English' },
    { key: 'zh-CN', name: '简体中文' }
  ]

  const menu = (
    <Menu style={{ borderRadius: 10, overflow: 'hidden' }}>
      <Menu.Item className='flex-column align-items-start' style={{ width: '100%', background: '#3370ff', display: 'flex', padding: 10 }}>
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
            {lang.key == app.language ? (
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
  const rlc = (
    <Menu style={{ borderRadius: 10, overflow: 'hidden' }}>
      <Menu.Item>{t('header.overview')}</Menu.Item>
      <Menu.Item>{t('header.quickstart')}</Menu.Item>
      <Menu.Item>{t('header.tip')}</Menu.Item>
      <Menu.Item>{t('header.contact')}</Menu.Item>
    </Menu>
  )
  return (
    <React.Fragment>
      <Header style={{ minWidth: '1000px' }} className='header'>
        <Skeleton loading={app.loading} active paragraph={{ rows: 0 }}>
          <div className='left'>
            <div style={{ cursor: 'pointer' }} onClick={menuClick}>
              {menuToggle ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              <span style={{ marginLeft: '6px' }}>{app.current}</span>
            </div>
          </div>
        </Skeleton>
        <div className='right'>
          <Skeleton loading={app.loading} active title={{ width: '100%' }} paragraph={{ rows: 0 }}>
            <div className='mr-3'>
              <Tooltip placement='bottom' title={t('back_to_home')}>
                <a
                  onClick={() => {
                    backHome()
                  }}
                >
                  <Avatar icon={<HomeOutlined />} alt='home' style={{ color: '#0677FF', backgroundColor:"transparent"}}/>
                </a>
              </Tooltip>
            </div>
            <div>
              <Dropdown placement='bottomRight' overlay={menu} overlayStyle={{ width: '20rem' }}>
                <div className='ant-dropdown-link'>
                  <Avatar src={user.avatar || avatar} alt='avatar' style={{ cursor: 'pointer' }} />
                </div>
              </Dropdown>
            </div>
          </Skeleton>
        </div>
      </Header>
    </React.Fragment>
  )
}

AppHeader.propTypes = {
  menuClick: PropTypes.func,
  menuToggle: PropTypes.bool,
  loginOut: PropTypes.func
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
export default connect(mapStateToProps, mapDispatchToProps)(AppHeader)
