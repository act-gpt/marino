import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Layout, Button, Empty } from 'antd'
import { useLocation, useHistory } from 'react-router-dom'
import { HomeOutlined, FileProtectOutlined } from '@ant-design/icons'
import '@/style/login.scss'

const Forbidden = (props) => {
  const home = '/'
  const { t, i18n } = useTranslation()
  return (
    <Layout className='login animated fadeIn'>
      <div className='model'>
        <div className='login-form'>
          <Empty
            image={<FileProtectOutlined />}
            imageStyle={{
              height: 200,
              lineHeight: '200px',
              fontSize: 100,
              color: '#FDB600'
            }}
            description={<p>{t('codes.403')}</p>}>
            {
              <Button
                onClick={() => {
                  window.location.replace(home)
                }}
                type='primary'
                size='large'
                shape='round'
              >
                {t('back_to_home')}
              </Button>
            }
          </Empty>
        </div>
      </div>
    </Layout>
  )
}

export default Forbidden
