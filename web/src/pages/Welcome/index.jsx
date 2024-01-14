import React, { useState, useEffect } from 'react'
import Api from '@/apiService/Api'
import { Layout, Input, Button, Divider, message, Form, Select } from 'antd'
import { useHistory } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { UserOutlined } from '@ant-design/icons'
import '@/style/login.scss'
import config from '@/config'
import 'react-phone-input-2/lib/style.css'

const msg = message

const Welcome = (props) => {
  const { t } = useTranslation()
  const history = useHistory()
  const user = JSON.parse(localStorage.getItem('user'))

  useEffect(() => {
    if (!user || user.org_id) {
      history.replace('/login')
    }
  })

  const [inputs, setInputs] = useState({
    name: '',
    contact: '',
    phone: '',
    size: '',
    roles: ''
  })

  const [loading, setLoading] = useState(false)

  function handleChange(e) {
    const { name, value } = e.target
    setInputs((inputs) => ({ ...inputs, [name]: value }))
  }

  async function handleSubmit(e) {
    if (!inputs.name) {
      msg.error(t('welcome.missing_org'))
      return
    }
    const conf = config.org.config
    conf.role = parseInt(inputs.role, 10)
    conf.size = parseInt(inputs.size, 10)
    const values = {
      ...inputs,
      ...config.org
    }
    setLoading(true)
    const res = await Api.add_org(values)
    const { success, message, data } = res
    if (success) {
      user.org_id = data.id
      localStorage.setItem('user', JSON.stringify(user))
      history.replace('/apps')
    } else {
      msg.error(message)
    }
    setLoading(false)
  }

  const roles = t('welcome.roles', { returnObjects: true })
  const sizes = t('welcome.sizes', { returnObjects: true })

  return (
    <Layout className='login animated fadeIn'>
      <div className='model'>
        <div className='login-form'>
          <h3 style={{ textAlign: 'center' }}>{t('welcome.org_info')}</h3>
          <Divider />
          <Form onFinish={handleSubmit}>
            <Form.Item
              name='name'
              rules={[
                {
                  required: true,
                  message: t('welcome.name_required')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='name' placeholder={t('welcome.name_placeholder')} />
            </Form.Item>
            <Form.Item
              name='contact'
              rules={[
                {
                  required: true,
                  message: t('welcome.contact_required')
                }
              ]}
            >
              <Input onChange={handleChange} prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='contact' placeholder={t('welcome.contact_placeholder')} />
            </Form.Item>
            <Form.Item
              name='phone'
              rules={[
                {
                  required: true,
                  message: t('welcome.phone_required')
                }
              ]}
            >
              {/*
                     <PhoneInput
                  country={"cn"}
                  masks={{"cn": "..........."}}
                  value={inputs.phone}
                  onChange={phone => setInputs((inputs) => ({ ...inputs, ['phone']: phone }))}
                  inputProps = {{
                      placeholder:t('welcome.phone_placeholder')
                    }
                  }
                />
                  */}

              {<Input onChange={handleChange} prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} name='phone' placeholder={t('welcome.phone_placeholder')} />}
            </Form.Item>
            <Form.Item
              name='size'
              rules={[
                {
                  required: true,
                  message: t('welcome.size_required')
                }
              ]}
            >
              <Select
                defaultValue=''
                onChange={(val) => {
                  setInputs((inputs) => ({ ...inputs, size: val }))
                }}
                style={{ width: '100%' }}
                options={sizes.map((v, i) => {
                  return {
                    value: i === 0 ? '' : i,
                    label: v
                  }
                })}
              />
            </Form.Item>
            <Form.Item
              name='roles'
              rules={[
                {
                  required: true,
                  message: t('welcome.role_required')
                }
              ]}
            >
              <Select
                defaultValue=''
                onChange={(val) => {
                  setInputs((inputs) => ({ ...inputs, roles: val }))
                }}
                style={{ width: '100%' }}
                options={roles.map((v, i) => {
                  return {
                    value: i === 0 ? '' : i,
                    label: v
                  }
                })}
              />
            </Form.Item>
            <Form.Item>
              <Button loading={loading} block type='primary' htmlType='submit'>
                {t('ok')}
              </Button>
            </Form.Item>
          </Form>
        </div>
      </div>
    </Layout>
  )
}

export default Welcome
