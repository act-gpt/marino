import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import {WarningOutlined, CheckOutlined, QuestionOutlined} from '@ant-design/icons'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import { Result,  List, Button, Badge, Card, Space, Avatar, notification } from 'antd'
import { AppContext } from "./index";

const Finished = (props) => {
  const { app, config , prev } = props
  const { t } = useTranslation()
  const [checked, setChecked] = useState({})
  const [disabled, setDisabled] = useState(true)
  const [index, setIndex] = useState(0)
  const [initialled, setInitialled] = useState({
    "ActGpt": false,
    "Baidu": false,
    "Db": false,
    "Embedding": false,
    "OpenAi": false
  })

  const map = {
    "db": "Db",
    "actgpt": "ActGpt",
    "baidu": "Baidu",
    "openai": "OpenAi",
    "embedding": "Embedding"
  }

  const data = app.language == "zh-CN" ? ["db", "actgpt", "baidu", "openai"] : ["db", "actgpt", "openai"]
  const services = t('setup.services', { returnObjects: true })
  
  
  const submit = async () => {
    if (config.Db.Dimension === 768){
      initialled["Embedding"] = true
    }
    setDisabled(true)
    initialled["Mail"] = config.Mail && config.Mail.SMTPToken ? true : false
    config.Initialled = initialled
    try{
      const res = await Api.save_config(config)
      const {code, data} = res
      if (!code){
        notification.success({
          message: t("setup.success_title"),
          description:  (<>{t("setup.success_save", { returnObjects: true }).map((val, idx)=> <div key={idx}>{val}</div>)}</>),
          duration: null,
          onClose: () => {
            setTimeout(()=>{
              window.location.replace("/login")
            }, 1000)
          }
        })
        return
      }
      setDisabled(false)
    }catch(e){
      if (e.response && e.response.status === 403){
        return notification.error({
          message: t("setup.message_title"),
          description:  t("codes." + e.response.status),
        })
      }
      setDisabled(false)
    }
    notification.error({
      message: t("setup.message_title"),
      description:  t("setup.success_failured"),
    })
  }
  
  const request = async (type, body) => {
    checked[type] = 0
    const idx = index + 1
    try{
      const res = await Api.check(type, body)
      if (!res.code){
        const key = map[type]
        checked[type] = 1
        initialled[key] = true
        setChecked({...checked})
        setDisabled(checked["db"] &&  modelPass() ?  false : true )
        setInitialled({...initialled})
      }else{
        checked[type] = 2
        setChecked({...checked})
      }
    }catch(e){
      checked[type] = 2
      setChecked({...checked})
    }
    setIndex(idx)
  }
  
  const modelPass = () => {
    return (["actgpt", "baidu", "openai"]).find((key)=> checked[key])
  }

  const next = (index) => {
    const type = data[index]
    if(type){
      check(type)
    }else{
      setDisabled(checked["db"] && modelPass() ?  false : true )
    }
  }

  const check = async (type) => {
    switch (type) {
      case "db":
        return request(type, config.Db);
      case "baidu":
        return request(type, config.Baidu);
      case "actgpt":
        return request(type, config.ActGpt);
      case "openai":
        return request(type, config.OpenAi);
      case "embedding":
        return request(type, config.Embedding);
      default:
        console.log(type)
    }
  }

  useEffect(()=>{
    setTimeout(() => next(index), 1000)
  },[index])

  useEffect(()=>{
    setTimeout(() => next(index), 1000)
  },[])

    return (<Card style={{ minWidth: "820px" }}>
      <Result
        status="info"
        title={t("setup.result_title")}
        extra={[
          <Space key="space" direction="vertical" size="middle" style={{ display: 'flex' }}>
            <List
              dataSource={data}
              renderItem={(item, idx) => (
                <List.Item key={idx}>
                  <div style={{paddingLeft: "80px"}}>
                    <span className="mr-3">{
                    checked[item] === 0 ? <span><span className="spin spin-small" style={{verticalAlign: "middle", width:"32px"}} /> </span> : 
                    checked[item] === 1 ?  <Avatar icon={<CheckOutlined />} style={{ backgroundColor: '#52c41a', color: '#fff' }} /> : 
                    checked[item] === 2 ? <Avatar icon={<WarningOutlined />} style={{ backgroundColor: '#fff', color: '#faad14' }} /> : 
                    <Avatar icon={<QuestionOutlined />}style={{ backgroundColor: '#0677ff', color: '#ffff' }} />} </span>
                    <span>{services[item]} {checked[item] === 2 && !(["db", "actgpt"]).find((key)=> key === item) ? t("setup.noworry") : ""} </span>
                  </div>
                </List.Item>
              )}/>
            <div className="mt-3"><Button size="large" className="mr-3" key="back" onClick={prev}>{t("setup.back")}</Button><Button disabled={disabled} onClick={submit} size="large" type="primary" key="save">{t("setup.save")}</Button></div>
          </Space>
        ]}
      />
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

export default connect(mapStateToProps, mapDispatchToProps)(Finished)