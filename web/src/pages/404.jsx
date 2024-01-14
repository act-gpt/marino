import React, { useState } from 'react'
import { Layout, Button } from 'antd'
import { useTranslation } from 'react-i18next'
import { useLocation, useHistory } from 'react-router-dom'
import '@/style/login.scss'

const NotFound = (props) => {
  const { t, i18n } = useTranslation()
  return (
    <Layout className='login animated fadeIn'>
      <div className='model'>
        <div className='login-form'>
          <h1>404</h1>
          <p className='py-5 text-center'>{t("codes.404")}</p>
          <div style={{ padding: '20px 0', textAlign: 'center' }}>
            <Button
              onClick={() => {
                window.location.replace('/')
              }}
              type='primary'
              shape='round'
            >
              {t('back')}
            </Button>
          </div>
        </div>
      </div>
    </Layout>
  )
}

export default NotFound
