import React, { useState } from 'react'
import Api from '@/apiService/Api'
import { useLocation, useHistory } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Layout, Input, Button, Divider, message, Form } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import '@/style/login.scss'
import util from '@/utils/chat'

const msg = message
const Login = (props) => {
  const { t } = useTranslation()
  const history = useHistory()
  const location = useLocation()
  const [loading, setLoading] = useState(false)

  /**
   * 登录
   * @param {*} e
   */
  const handleSubmit = (e) => {
    setLoading(true)
    Api.login(e)
      .then((res) => {
        setLoading(false)
        const { success, data, message } = res
        if (success) {
          localStorage.setItem('user', JSON.stringify(data))
          const redirect = util.params(location.search, 'redirect')
          if (!data.org_id) {
            return history.replace('/welcome')
          }
          if (redirect) {
            return history.replace(redirect)
          }
          history.replace('/apps')
        } else {
          msg.warning(message)
        }
      })
      .catch((err) => {
        msg.warning(err.message)
        setLoading(false)
      })
  }
  return (
    <Layout className='login animated fadeIn'>
      <div className='model'>
        <div className='login-form'>
          <h3 style={{ textAlign: 'center' }}>{t('login.info')}</h3>
          <Divider />
          <Form onFinish={handleSubmit}>
            <Form.Item
              name='username'
              rules={[
                {
                  required: true,
                  message: t('login.name_required')
                }
              ]}
            >
              <Input prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} placeholder={t('login.name_placeholder')} />
            </Form.Item>
            <Form.Item
              name='password'
              rules={[
                {
                  required: true,
                  message: t('login.password_required')
                }
              ]}
            >
              <Input prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} type='password' placeholder={t('login.password_placeholder')} />
            </Form.Item>
            <Form.Item>
              <Button loading={loading} block type='primary' htmlType='submit'>
                {t('continue')}
              </Button>
            </Form.Item>
          </Form>
          <p dangerouslySetInnerHTML={{ __html: t('login.signin') }}></p>
        </div>
      </div>
    </Layout>
  )
}

export default Login
