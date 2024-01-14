import React, { useState } from 'react'
import { message, Form, Input, Button } from 'antd'
import { useTranslation } from 'react-i18next'

import { useHistory } from 'react-router-dom'
import { connect } from 'react-redux'
import actions from '@/redux/actions'
import config from '@/config'
import Api from '@/apiService/Api'

import CommonModal from '@/commonComponents/commonModal'

const Create = (props) => {
  const history = useHistory()

  const [form] = Form.useForm()
  const { navigation, user, app, setHelpdesk, setApp, visible, setVisible } = props
  const [disabled, setDisabled] = useState(false)
  const { t, i18n } = useTranslation()

  const hide = () => {
    form.resetFields()
    setVisible(false)
  }

  const submit = () => {
    setDisabled(true)
    form.validateFields().then(async (values) => {
      if (!values.name) {
        message.error('Please enter bot name')
        return
      }
      const item = {
        name: values.name,
        ...config.bot
      }
      const res = await Api.add_bot(item)
      const { success, data } = res
      if (success) {
        message.success(res.message)
        props.setBot(data)
        history.push('/admin/' + data.id + '/knowledges')
        return
      } else {
        message.warn(res.message)
      }
    })
    setDisabled(false)
    //hide()
  }
  const rules = [{ required: true, message: t('required') }]
  return (
    <CommonModal
      title={t('left.create_title')}
      width={600}
      footer={false}
      visible={visible}
      common_cancel={hide}
      children={
        <Form layout='vertical' requiredMark={false} form={form} onFinish={submit}>
          <div className='py-3'>{t('left.create_desc')}</div>
          <Form.Item rules={rules} name='name' label={t('left.create_tip')}>
            <Input />
          </Form.Item>
          <Form.Item className='pt-5'>
            <div className='d-flex flex-row-reverse'>
              <Button type='primary' htmlType='submit' disabled={disabled}>
                {' '}
                {t('ok')}{' '}
              </Button>
              <Button className='mr-5' onClick={hide}>
                {' '}
                {t('cancel')}{' '}
              </Button>
            </div>
          </Form.Item>
        </Form>
      }
    />
  )
}

const mapStateToProps = (state) => {
  return {
    user: state.user,
    app: state.app,
    servicedesk: state.servicedesk
  }
}
const mapDispatchToProps = {
  ...actions
}
export default connect(mapStateToProps, mapDispatchToProps)(Create)
