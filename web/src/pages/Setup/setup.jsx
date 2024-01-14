import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import { Button, Form, Input, Space, Badge, Card, Switch, Alert, Radio, Divider, message, notification } from 'antd'

const Setup = (props) => {
  const { app, config, next, save } = props
  const [form] = Form.useForm()
  const { t } = useTranslation()
  const [more, setMore] = useState(false)
  const [checked, setChecked] = useState({})
  const [openaiVisible, setOpenaiVisible] = useState(false);

  const baiduUrl = "https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application"
  const url = process.env.REACT_CONSOLE_URL || "https://act-gpt.com"

  message.config({top: 60})

  const submit = async (vals) => {
    if (vals.Db.Dimension == 768){
      delete vals.Embedding
    }
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

  const validator = (item, value, cb, source, options) => {
    const s = checked[item.field]
    if (value && s == 1){
      return Promise.resolve()
    }
    return Promise.reject(new Error(""));
  }
  
  useEffect(() => {
    form.setFieldsValue(config)
    setMore(config.Db.Dimension==1536)
   },[config])

  return (
        <Card style={{ minWidth: "820px" }}>
          <Alert banner message={t('setup.tips', { returnObjects: true , url: url, baidu: baiduUrl}).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" className="mt-2" />
          <Form form={form} layout="vertical" scrollToFirstError={true} initialValues={config} autoComplete="off" onFinish={submit}>
            <Divider orientation="left">{t("setup.org")}</Divider>
            <Form.Item name={["Organization", "Name"]} label={t("setup.org_name")} rules={[{ required: true }]}>
              <Input placeholder={t("setup.org_name")} count={{show: true, max: 18,}}/>
            </Form.Item>
            <Divider orientation="left">{t("setup.address")}</Divider>
            <Form.Item name={["Db", "DataSource"]} label={t("setup.address")} rules={[{ required: true }]} help={t("setup.address_help")}  hasFeedback validateStatus={checked["db"] == 1 ? "success" : checked["db"] == 2 ? "error" : ""}>
              <Space.Compact block>
                <Input placeholder={t("setup.address")} defaultValue={config.Db.DataSource}/>
                <Button style={{ width: 80 }} onClick={() => { check('db')}}>{t("setup.test")}</Button>
              </Space.Compact>
            </Form.Item>
            <Form.Item name={["Db", "Redis"]} label={t("setup.redis")} help={t("setup.redis_help")}>
              <Input />
            </Form.Item>
            <Form.Item name={["Db", "Dimension"]} label={t("setup.dimension")} help={t("setup.dimension_help")}>
              <Radio.Group buttonStyle="solid" onChange={ ({target: { value }} )=> setMore(value==1536)}>
                <Radio.Button value={768}>ACT GPT</Radio.Button>
                <Radio.Button value={1536}>Open AI</Radio.Button>
              </Radio.Group>
            </Form.Item>
            <Alert banner message={t('setup.banner', { returnObjects: true, url: url }).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" />
            {
              more ? (<>
                <Form.Item name={["Embedding", "Host"]} label={t("setup.emd_host")}>
                  <Input placeholder={t("setup.emd_host")} />
              </Form.Item>
              <Form.Item name={["Embedding", "Api"]} label={t("setup.emd_api")}>
                <Input placeholder={t("setup.emd_api")} />
              </Form.Item>
              <Form.Item name={["Embedding", "Model"]} label={t("setup.emd_model")}>
                <Input placeholder={t("setup.emd_model")} />
              </Form.Item>
              <Form.Item name={["Embedding", "AccessKey"]} label={t("setup.emd_key")} hasFeedback validateStatus={checked["embedding"] == 1 ? "success" : checked["embedding"] == 2 ? "error" : ""}>
                <Space.Compact block>
                  <Input.Password
                    placeholder={t("setup.emd_key")}
                    visibilityToggle={{ visible: openaiVisible, onVisibleChange: setOpenaiVisible }} />
                  <Button style={{ width: 80 }} onClick={() => { check('embedding')}}>{t("setup.test")}</Button>
                </Space.Compact>
              </Form.Item>
              </>) : ""
            }
            {
              app.language == "zh-CN" ? <>
              <Divider orientation="left">{t("setup.moderation")}</Divider>
              <Form.Item name={["Moderation", "CheckContent"]} label={t("setup.mdt")} help={t("setup.mdt_help")}>
                <Switch checkedChildren={t("setup.checked")} unCheckedChildren={t("setup.unchecked")} defaultChecked={false} />
              </Form.Item>
              <Form.Item name={["Moderation", "Api"]} label={t("setup.mtda")}>
                <Input placeholder={t("setup.mtda")} />
              </Form.Item>
              </> : ""
            }
            {
              process.env.MAIL_SERVICE === "true" ? (<><Divider orientation="left">{t("setup.email")}</Divider>
              <Form.Item name={["Mail", "SMTPToken"]} label={t("setup.resend_token")}  help={t("setup.resend_token_help")}>
                <Input.Password
                      placeholder={t("setup.resend_token")}/>
              </Form.Item>
              <Form.Item name={["Mail", "SMTPFrom"]} label={t("setup.resend_form")}>
                <Input placeholder={t("setup.resend_form")} />
              </Form.Item>
              </>) : ""
            }
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

export default connect(mapStateToProps, mapDispatchToProps)(Setup)