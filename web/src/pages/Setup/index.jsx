import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import { Image, Layout,theme, FloatButton, Steps, Badge, Button, Switch } from 'antd'
import { CustomerServiceOutlined, CommentOutlined} from '@ant-design/icons'
import Contact from '@/components/Contact'
import Normal from './normal'
import Setup from './setup'
import Model from './model'
import Finished from './finished'
import Api from '@/apiService/Api'
import logo from '@/assets/images/logo.png'

const { Header, Content, Footer } = Layout;
export const AppContext = React.createContext({})

const Index = (props) => {
  const { app, user, config, setConfig } = props
 
  const { t } = useTranslation()
  const ref = useRef()
  const { token } = theme.useToken();
  const contentStyle = {
    lineHeight: '260px',
    color: token.colorTextTertiary,
    backgroundColor: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: `1px dashed ${token.colorBorder}`,
    marginTop: 36,
  };

  const [current, setCurrent] = useState(0)
  const [advanced, setAdvanced] = useState(false)
  const [visibleContact, setVisibleContact] = useState(false)

  const next = (setp) => {
    setCurrent(current + 1)
    scroll()
  };

  const prev = (setp) => {
    setCurrent(current - 1)
    scroll()
  };

  const save = (conf) => {
    setConfig({...conf})
  }

  const scroll = () => {
    ref.current?.scrollIntoView('auto')
  }

  const setModel = (checked) => {
    setCurrent(0)
    setAdvanced(checked)
  }

  const advanceds = [
    {
      title: t("setup.step1"),
      content: <Setup  next={next} save={save} />,
    },
    {
      title: t("setup.step2"),
      content:  <Model prev={prev} next={next} save={save}/>,
    },
    {
      title: t("setup.step3"),
      content: <Finished prev={prev}/>,
    }
  ]
  const normal = [
    {
      title: t("setup.step4"),
      content: <Normal  next={next} save={save} />,
    },
    {
      title: t("setup.step3"),
      content: <Finished prev={prev}/>,
    }
  ]

  const normal_items = normal.map((item) => ({ key: item.title, title: item.title }))
  const advanceds_items = advanceds.map((item) => ({ key: item.title, title: item.title }))

  useEffect(() => {
    Api.config().then((res)  => {
      const {Initialled} = res.data
      if (Initialled.Db) {
        return window.location.replace("/login")
      }
    })
   }, [])

   useEffect(() => {
    console.log(current)
  }, [current])
  console.log("render", current)

  return (<Layout className='app'>
    <Header
      className='header justify-content-center'>
      <div className='d-flex align-items-center justify-content-between' style={{ minWidth: "820px"}}>
        <div className='d-flex align-items-center'>
          <Image
            className="mr-3"
            style={{ width: '28px', height: '28px' }}
            src={logo}
            preview={false} /> <h3>Marino</h3>
        </div>
        <Switch onChange={setModel} checked={advanced} checkedChildren={t("setup.advanced")} unCheckedChildren={t("setup.normal")} />
      </div>
    </Header>
    {
      <Content className='d-flex justify-content-center' style={{
      minHeight: "calc(100vh - 100px)", padding: 24,
      background: token.colorBgContainer,
      borderRadius: token.borderRadiusLG}}>
        <div style={{ minWidth: "820px", padding: "32px 0 0"}} ref={ref} >
          {
          advanced ?
          (<><Steps current={current} items={advanceds_items} />
          <div style={contentStyle}>
            <Badge.Ribbon text={process?.env?.REACT_APP_VERSION || "0.1.0"}>
            {advanceds[current]?.content}
            </Badge.Ribbon>
          </div></>) : (<><Steps current={current} items={normal_items} />
            <div style={contentStyle}>
              <Badge.Ribbon text={process?.env?.REACT_APP_VERSION || "0.1.0"}>
                {normal[current]?.content}
              </Badge.Ribbon>
            </div></>)
          }
        </div>
    </Content>
    }
    
    <Footer style={{ textAlign: 'center', height: "50px" }}>{config.Organization.Name}</Footer>
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
        {app.language == "zh-CN" 
        ? <FloatButton  tooltip={<div>{t('home.contact_us')}</div>} icon={<CommentOutlined />} onClick={()=> setVisibleContact(true)} />
        : <FloatButton  tooltip={<div>{t('discord')}</div>} href={"https://discord.gg/GJrmRDyh"} target="_blank"/>}
      </FloatButton.Group>
  </Layout>)
}
const mapStateToProps = (state) => {
  return {
    app: state.app,
    user: state.user,
    config: state.config
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Index)