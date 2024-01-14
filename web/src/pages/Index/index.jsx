import React, { useState, useEffect } from 'react'
import { Image, Layout, Button, Result, Row , Col } from 'antd'
import { useTranslation } from 'react-i18next'
import { SmileOutlined } from '@ant-design/icons'
import logo from '@/assets/images/logo.png'

const { Header, Content, Footer } = Layout

const Index = (props) => {
  const { t, i18n } = useTranslation()

  window.addEventListener('beforeinstallprompt', (event) => {
    // Prevent the mini-infobar from appearing on mobile.
    event.preventDefault();
    const host = process.env.REACT_APP_WEB || 'act-gpt.com'
    if(localStorage.getItem('install.alert') || !(new RegExp(host)).test(window.location.href)){
      return
    }
    window.deferredPrompt = event
  })
  
  return (
    <Layout className='home vh-100'>
      <Header className='header py-2' style={{ background: "#fff"}}>
          <Row>
            <Col flex='1 1 200px'>
            <div className='d-flex align-items-center'>
            <Image
              className="mr-3"
              style={{ width: '28px', height: '28px' }}
              src={logo}
              preview={false}/>
              <h3>Marino</h3>
            </div>
            </Col>
            <Col flex='0 1 60px'>
              
            </Col>
          </Row>
        </Header>
      <Layout className='h-100' style={{ padding: 0, background: '#fff' }}>
        <Content>
          <div className='h-100 text-center d-flex justify-content-center align-items-center'>
            <Result
              icon={<SmileOutlined />}
              title={t('index.intro')}
              subTitle={t('index.intro_subtitle')}
              extra={
                <Button
                  type='primary'
                  size='large'
                  shape="round"
                  onClick={() => {
                    props.history.push(`/apps`)
                  }}>
                  {t('index.intro_button')}
                </Button>
              }
            />
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default Index
