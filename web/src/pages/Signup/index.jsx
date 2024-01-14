import React, { useState } from 'react'
import Api from '@/apiService/Api'
import { useHistory } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Layout, Input, Button, Divider, message, Form, Row, Col } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined, FieldBinaryOutlined } from '@ant-design/icons'
import '@/style/login.scss'

const msg = message

const Signup = (props) => {
  const { t } = useTranslation()
  const history = useHistory()

  const [inputs, setInputs] = useState({
    username: '',
    password: '',
    password2: '',
    email: '',
    verification_code: ''
  })
  const { username, password, password2 } = inputs
  const [loading, setLoading] = useState(false)
  const [send, setSend] = useState(false)

  function handleChange(e) {
    const { name, value } = e.target
    setInputs((inputs) => ({ ...inputs, [name]: value }))
  }

  async function handleSubmit(e) {
    if (password.length < 8) {
      msg.warning(t('signin.password_error'))
      return
    }
    if (password !== password2) {
      msg.warning(t('signin.comfirm_error'))
      return
    }
    if (username && password) {
      setLoading(true)
      const res = await Api.signup(inputs)
      const { success } = res
      if (success) {
        history.push('/login')
        msg.success(t('signin.success'))
      } else {
        msg.error(t('signin.failure'))
      }
      setLoading(false)
    }
  }

  const sendVerificationCode = async () => {
    if (inputs.email === '') return
    setSend(true)
    const res = await Api.verification(inputs.email)
    const { success, message } = res
    if (success) {
      msg.success(t('signin.send_success'))
    } else {
      msg.error(message)
    }
    setTimeout(() => {
      setSend(false)
    }, 1000 * 60)
  }

  return (
    <Layout className='login animated fadeIn'>
      <div className='model'>
        <div className='login-form'>
          <h3 style={{ textAlign: 'center' }}>{t('signin.info')}</h3>
          <Divider />
          <Form onFinish={handleSubmit}>
            <Form.Item
              name='username'
              rules={[
                {
                  required: true,
                  message: t('signin.name_required')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='username' placeholder={t('signin.name_placeholder')} />
            </Form.Item>
            <Form.Item
              name='password'
              rules={[
                {
                  required: true,
                  message: t('signin.password_required')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='password' type='password' placeholder={t('signin.password_placeholder')} />
            </Form.Item>
            <Form.Item
              name='password2'
              rules={[
                {
                  required: true,
                  message: t('signin.password_required')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='password2' type='password' placeholder={t('signin.comfirm_password_placeholder')} />
            </Form.Item>

            <Form.Item noStyle>
              <Row gutter={8}>
                <Col span={18}>
                  <Form.Item
                    name='email'
                    rules={[
                      {
                        required: true,
                        message: t('signin.email_required')
                      }
                    ]}
                  >
                    <Input onChange={handleChange} prefix={<MailOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='email' type='email' placeholder={t('signin.email_placeholder')} />
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Button onClick={sendVerificationCode} disabled={send}>
                    {t('signin.send')}
                  </Button>
                </Col>
              </Row>
            </Form.Item>

            <Form.Item
              name='verification_code'
              rules={[
                {
                  required: true,
                  message: t('signin.captcha_placeholder')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<FieldBinaryOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='verification_code' type='text' placeholder={t('signin.captcha_placeholder')} />
            </Form.Item>
            <Form.Item>
              <Button loading={loading} block type='primary' htmlType='submit'>
                {t('continue')}
              </Button>
            </Form.Item>
          </Form>
          <p dangerouslySetInnerHTML={{ __html: t('signin.login') }}></p>
        </div>
      </div>
    </Layout>
  )
}

export default Signup
