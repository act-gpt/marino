import React, { useState, useReducer, useEffect } from 'react'

import { Route, Switch, Redirect, withRouter } from 'react-router-dom'
import { Layout, BackTop, message, Skeleton } from 'antd'

import router from '@/router'
import init from '@/utils/init'
// import echarts from 'echarts/lib/echarts'
import menus from './menu'
import '@/style/layout.scss'

import AppHeader from './AppHeader.jsx'
import AppAside from './AppAside.jsx'

import { connect } from 'react-redux'
import actions from '@/redux/actions'
import Forbidden from '@/pages/403'
import Loading from '@/components/Loading'
import NotFound from '@/pages/404'
import Api from '../apiService/Api'
import { useLocation, useHistory } from 'react-router-dom'

const { Content } = Layout
//const history = useHistory()
const getMenu = (menu) => {
  return menu.filter((res) => res.auth && res.auth.indexOf(1) !== -1).slice(0)
}

class DefaultLayout extends React.Component {
  constructor(props) {
    super(props)
    this.props = props
    this.loading = true
    this.state = {
      menu: getMenu(menus),
      auth: true
    }
    this.init()
  }
  init() {    
    init(this.props).then(async(res) => {
      this.props.setApp({ loading: false })
      this.loading = false
      if (!res.permission) {
        this.setState({ auth: false })
        await Api.signout()
        this.props.history.replace('/login?redirect=' + decodeURIComponent(window.location.pathname))
      }
    }).catch(()=> this.loading = false)
  }

  menuClick() {
    this.props.tiggerToggle()
  }

  async loginOut() {
    localStorage.removeItem('user')
    await Api.signout()
    this.props.history.push('/login')
  }

  render() {
    const { app, user, toggle } = this.props
    const { menu, auth } = this.state
    return !auth ? (
      <Forbidden />
    ) : (
      <Layout className='app'>
        <AppAside menuToggle={toggle.menuToggle} menu={menu ? menu : []} navigation={this.props} />
        <Layout
          style={{
            marginLeft: toggle.menuToggle ? '80px' : '250px',
            minHeight: '100vh',
            transition: 'min-width .2s',
            overflowX: 'auto'
          }}
        >
          <AppHeader menuToggle={toggle.menuToggle} menuClick={this.menuClick.bind(this)} loginOut={this.loginOut.bind(this)} />
          <Content style={{ minWidth: '1000px' }} className='content'> 
           {this.loading ? <Loading /> : (
              <Skeleton loading={app.loading} active title={false} paragraph={{ rows: 6 }} style={{ padding: '20px' }}>
              <Switch>
                {router.map((item) => {
                  return <Route key={item.path} path={item.path} exact={item.exact} render={(props) => <item.component {...props} />}></Route>
                })}
                <Route
                  key={'/'}
                  path={'/'}
                  exact={false}
                  render={(props) => {
                    props.history.push('/')
                    return <h1>Home</h1>
                  }}
                ></Route>
              </Switch>
            </Skeleton>
           )}
          </Content>
        </Layout>
      </Layout>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    toggle: state.toggle,
    user: state.user,
    app: state.app
  }
}

const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(DefaultLayout))
