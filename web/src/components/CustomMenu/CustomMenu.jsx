import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { Menu, Skeleton } from 'antd'
import { Link, withRouter } from 'react-router-dom'
import { createFromIconfontCN } from '@ant-design/icons'

import { connect } from 'react-redux'
import actions from '@/redux/actions'

// 处理 pathname
const getOpenKeys = (string) => {
  let newStr = '',
    newArr = [],
    arr = string.split('/').map((i) => '/' + i)
  for (let i = 1; i < arr.length - 1; i++) {
    newStr += arr[i]
    newArr.push(newStr)
  }
  return newArr
}

const IconFont = createFromIconfontCN({
  scriptUrl: '/js/icon.js'
})

const CustomMenu = (props) => {
  const { user, app, bot, setApp } = props
  const [state, setstate] = useState({
    openKeys: [],
    selectedKeys: []
  })

  let { openKeys, selectedKeys } = state

  const menu = (props.menu || []).slice(0)

  const replace = (path) => path.replace(':bot', bot.id || 0)

  // 页面刷新的时候可以定位到 menu 显示
  useEffect(() => {
    let { pathname } = props.location
    let f = menu.find((item) => pathname.startsWith(replace(item.key)))
    const current = (f && f.title) || ''
    if (current && current != app.current) {
      setApp({ current: (f && f.title) || '' })
    }
    setstate((prevState) => {
      return {
        ...prevState,
        selectedKeys: [replace((f && f.key) || '')],
        openKeys: getOpenKeys(replace((f && f.key) || ''))
      }
    })
  }, [bot])

  // 只展开一个 SubMenu
  const onOpenChange = (openKeys) => {
    setstate((prevState) => {
      if (openKeys.length === 0 || openKeys.length === 1) {
        return { ...prevState, openKeys }
      }
      const latestOpenKey = openKeys[openKeys.length - 1]

      // 这里与定义的路由规则有关
      if (latestOpenKey.includes(openKeys[0])) {
        return { ...prevState, openKeys }
      } else {
        return { ...prevState, openKeys: [latestOpenKey] }
      }
    })
  }

  const renderMenuItem = ({ key, icon, title }) => {
    key = replace(key)
    return (
      <Menu.Item key={key}>
        <Skeleton loading={app.loading} active avatar paragraph={{ rows: 0 }}>
          <Link
            to={(location) => {
              return { ...location, pathname: key }
            }}
            onClick={(e) => {
              setApp({ current: title })
            }} >
            <IconFont type={icon} style={{ fontSize: '1.0rem' }} />
            <span style={{ fontSize: '1.0rem', textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>{title}</span>
          </Link>
        </Skeleton>
      </Menu.Item>
    )
  }

  // 循环遍历数组中的子项 subs ，生成子级 menu
  const renderSubMenu = ({ key, icon, title, subs }) => {
    return (
      <Menu.SubMenu
        key={key}
        title={
          <span>
            <IconFont type={icon} />
            <span style={{ display: 'inline-block' }}>{title}</span>
          </span>
        }
      >
        {subs &&
          subs.map((item) => {
            return item.subs && item.subs.length > 0 ? renderSubMenu(item) : renderMenuItem(item)
          })}
      </Menu.SubMenu>
    )
  }

  return (
    <Menu
      mode='inline'
      openKeys={openKeys}
      selectedKeys={props.menuToggle === true ? '' : selectedKeys}
      onClick={({ key }) => setstate((prevState) => ({ ...prevState, selectedKeys: [key] }))}
      onOpenChange={onOpenChange}>
      {menu.map((item) => {
        return item.subs && item.subs.length > 0 ? renderSubMenu(item) : renderMenuItem(item)
      })}
    </Menu>
  )
}

CustomMenu.propTypes = {
  menu: PropTypes.array.isRequired
}

const mapStateToProps = (state) => {
  return {
    app: state.app,
    user: state.user,
    bot: state.bot
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(CustomMenu))
