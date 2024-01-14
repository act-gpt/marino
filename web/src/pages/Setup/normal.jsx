import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import { AppContext } from "./index";

import { Button, Form, Input, Space, Badge, Card, Switch, Alert, Radio, Divider, message, notification } from 'antd'

const Model = (props) => {
  const { app, config, save, next } = props

  const [form] = Form.useForm()
  const { t } = useTranslation()
  const [more, setMore] = useState(false)
  const [checked, setChecked] = useState({})
  const [baiduVisible, setBaiduVisible] = useState(false);
  const [actgptVisible, setActgptVisible] = useState(false);
  const [openaiVisible, setOpenaiVisible] = useState(false);

  const baiduUrl = "https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application"
  const url = process.env.REACT_CONSOLE_URL || "https://act-gpt.com"
  const alertTip = (<>{t('setting.alert_tip', { returnObjects: true }).map((val, idx) => <div key={idx}>{val}</div>)}</>)
 
  message.config({top: 60})

  const submit = async (vals) => {
    save({...vals})
    next()
  }
  const request = async (type, body) => {
    const update = {...checked}
    message.loading({
      content: <div style={{width: "180px"}}>{t("setup.loading")}</div>,
      duration: 0,
      key: type
    })
    try{
      const res = await Api.check(type, body)
      if (!res.code){
        update[type] = 1
        setChecked(update)
        message.destroy(type)
        return notification.success({
          message: t("setup.message_title"),
          description:  t("setup.success_desc"),
        })
      }
    }catch(e){
      if (e.response.status === 429){
        return notification.warning({
          message:  t("setup.message_title"),
          description: t("setup.message_retry"),
        })
      }
    }
    message.destroy(type)
    update[type] = 2
    setChecked(update)
    notification.error({
      message:  t("setup.message_title"),
      description:  t("setup.failure_desc"),
    })
  }

  const check = async (type) => {
    const items = form.getFieldsValue()
    switch (type) {
      case "db":
        return request(type, items.Db);
      case "baidu":
        return request(type, items.Baidu);
      case "actgpt":
        return request(type, items.ActGpt);
      case "openai":
        return request(type, items.OpenAi);
      case "embedding":
        return request(type, items.Embedding);
      default:
        console.log(type)
    }
  }
  const validator = (item, value) => {
    const s = checked[item.field]
    if (value && s == 1){
      return Promise.resolve()
    }
    return Promise.reject(new Error(""));
  }
  return (
        <Card style={{ minWidth: "820px" }}>
          <Alert banner message={t('setup.normal_tips', { returnObjects: true , url: url, baidu: baiduUrl}).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" className="mt-2" />
          <Form form={form} layout="vertical" initialValues={config} autoComplete="off" onFinish={submit}>
          <Divider orientation="left">{t("setup.org")}</Divider>
            <Form.Item name={["Organization", "Name"]} label={t("setup.org_name")} rules={[{ required: true }]}>
              <Input placeholder={t("setup.org_name")} count={{show: true, max: 18,}}/>
            </Form.Item>
            <Divider orientation="left">{t("setup.address")}</Divider>
            <Form.Item name={["Db", "DataSource"]} label={t("setup.address")} rules={[{ required: true}]} help={t("setup.address_help")}>
              <Space.Compact block>
                <Input placeholder={t("setup.address")} defaultValue={config.Db.DataSource}/>
                <Button style={{ width: 80 }} onClick={() => { check('db')}}>{t("setup.test")}</Button>
              </Space.Compact>
            </Form.Item>
            <Form.Item name={["Db", "Dimension"]} label={t("setup.dimension")} help={t("setup.dimension_help")}>
              <Radio.Group buttonStyle="solid" onChange={ ({target: { value }} )=> setMore(value==1536)}>
                <Radio.Button value={768}>ACT GPT</Radio.Button>
                <Radio.Button value={1536} disabled>Open AI</Radio.Button>
              </Radio.Group>
            </Form.Item>
            <Divider orientation="left">{t("setup.actgpt")}</Divider>
            <Form.Item name={["ActGpt", "AccessKey"]} label={t("setup.actgpt_key")} rules={[{ required: true }]}>
              <Space.Compact block>
                <Input.Password
                  placeholder={t("setup.actgpt_key")}
                  visibilityToggle={{ visible: actgptVisible, onVisibleChange: setActgptVisible }} />
                <Button style={{ width: 80 }} onClick={() => { check('actgpt')}}>{t("setup.test")}</Button>
              </Space.Compact>
            </Form.Item>
            <Form.Item className="pt-4">
              <Button type='primary' size="large" htmlType='submit'>
                {t('continue')}
              </Button>
            </Form.Item>
          </Form>
        </Card>)
}

const mapStateToProps = (state) => {
  return {
    app: state.app,
    config: state.config
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Model)