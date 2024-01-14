import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import { Button, Form, Input, Space, Badge, Card, Switch, Alert, Radio, Divider, message, notification } from 'antd'

const Model = (props) => {
  const { app, config, prev, next, save,} = props

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
  return (
        <Card style={{ minWidth: "820px" }}>
          {
              app.language == "zh-CN" ? <Alert className="mt-2" banner message={alertTip} type="warning" /> : ""
          }
          <Form form={form} layout="vertical" initialValues={config} autoComplete="off" onFinish={submit}>
            <Divider orientation="left">{t("setup.actgpt")}</Divider>
            <Form.Item name={["ActGpt", "AccessKey"]} label={t("setup.actgpt_key")} hasFeedback validateStatus={checked["actgpt"] == 1 ? "success" : checked["actgpt"] == 2 ? "error" : ""}>
              <Space.Compact block>
                <Input.Password
                  placeholder={t("setup.actgpt_key")}
                  visibilityToggle={{ visible: actgptVisible, onVisibleChange: setActgptVisible }} />
                <Button style={{ width: 80 }} onClick={() => { check('actgpt')}}>{t("setup.test")}</Button>
              </Space.Compact>
            </Form.Item>
            <Alert banner message={t('setup.actgpt_alert', { returnObjects: true, url: url },).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" />
            {
              app.language != "zh-CN" ? "" : (<> <Divider orientation="left">{t("setup.baidu")}</Divider>
                <Form.Item name={["Baidu", "ClientId"]} label={t("setup.baidu_key")} hasFeedback  validateStatus={checked["baidu"] == 1 ? "success" : checked["baidu"] == 2 ? "warning" : ""}>
                  <Input placeholder={t("setup.baidu_key")} />
                </Form.Item>
                <Form.Item name={["Baidu", "ClientSecret"]} label={t("setup.baidu_secret")} hasFeedback  validateStatus={checked["baidu"] == 1 ? "success" : checked["baidu"] == 2 ? "warning" : ""}>
                  <Space.Compact block>
                    <Input.Password
                      placeholder={t("setup.baidu_secret")}
                      visibilityToggle={{ visible: baiduVisible, onVisibleChange: setBaiduVisible }}
                    />
                    <Button style={{ width: 80 }} onClick={() => { check('baidu') }}>{t("setup.test")}</Button>
                  </Space.Compact>
                </Form.Item>
                <Alert banner message={t('setup.baidu_alert', { returnObjects: true , url: baiduUrl}).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" /></>)
            }
            <Divider>{t("setup.openai")}</Divider>
            <Form.Item name={["OpenAi", "Type"]} label={t("setup.openai_type")}>
              <Radio.Group buttonStyle="solid">
                <Radio.Button value={"openai"}>Open AI</Radio.Button>
                <Radio.Button value={"azure"}>Azure</Radio.Button>
              </Radio.Group>
            </Form.Item>
            <Form.Item name={["OpenAi", "AccessKey"]} label={t("setup.openai_key")} hasFeedback validateStatus={checked["openai"] == 1 ? "success" : checked["openai"] == 2 ? "warning" : ""}>
              <Space.Compact block>
                <Input.Password
                  placeholder={t("setup.openai_key")}
                  visibilityToggle={{ visible: openaiVisible, onVisibleChange: setOpenaiVisible }} />
                <Button style={{ width: 80 }} onClick={() => { check('openai')}}>{t("setup.test")}</Button>
              </Space.Compact>
            </Form.Item>
            <Form.Item name={["OpenAi", "Host"]} label={t("setup.openai_host")} help={t("setup.openai_host_help")}>
              <Input placeholder={t("setup.openai_host")} />
            </Form.Item>
            <Form.Item name={["OpenAi", "APIVersion"]} label={t("setup.openai_v")}  help={t("setup.openai_v_help")} >
              <Input placeholder={t("setup.openai_v_hp")} />
            </Form.Item>
            {app.language != "zh-CN" ? "" :
              <Alert banner message={t('setup.openal_alert', { returnObjects: true }).map((val, idx)=><div key={idx} dangerouslySetInnerHTML={{ __html: val }}/>)} type="info" />}
            <Form.Item className="pt-4">
              <Button size="large" className="mr-3" onClick={prev}>
                {t("setup.back")}
              </Button>
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