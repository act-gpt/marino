import React, { useState, useRef, useEffect} from 'react'
import { Avatar, Form, Input, Button, Alert, Row, Col, Slider, Select, Tooltip, message , Tour} from 'antd'
import { QuestionCircleOutlined } from '@ant-design/icons'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import Contact from '@/components/Contact'
import './style.scss'

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 24 }
}

const Setting = (props) => {
  const { bot, app } = props
  const [form] = Form.useForm()
  const { t } = useTranslation()

  const [open, setOpen] = useState(false)
  const [setting, setSetting] = useState(bot.setting)
  const [visible, setVisible] = useState(false)
  const [disable, setDisable] = useState(false)
  const [models, setModels] = useState([])
  const ref1 = useRef(null);
  const ref2 = useRef(null);
  const ref3 = useRef(null);
  const ref4 = useRef(null);

  const steps = [
    {
      title: t("setting.tour.ref1_title"),
      description: t("setting.tour.ref1_desc"),
      target: () => ref1.current,
      placement:"bottom",
    },
    {
      title: t("setting.tour.ref2_title"),
      description: t("setting.tour.ref2_desc"),
      target: () => ref2.current,
      placement:"top",
    },
    {
      title: t("setting.tour.ref3_title"),
      description: t("setting.tour.ref3_desc"),
      target: () => ref3.current,
      placement:"top",
    },
    {
      title: t("setting.tour.ref4_title"),
      description: t("setting.tour.ref4_desc"),
      target: () => ref4.current,
      placement:"top",
    },
  ]
  
  useEffect(() => {
    if (!localStorage.getItem('setting.tour') && ref1.current){
      setOpen(true)
    }
    getModels()
  }, [])
  
  const getModels = async () => {
    if(models.length){
      return
    }
    let init = {
      "baidu": false,
      "actgpt": false,
      "openai": false,
    }
    try{
      const res = await Api.config()
      const {Initialled} = res.data
      Object.keys(Initialled).forEach((key)=>{
        init[key.toLowerCase()] = Initialled[key]
      })
    }catch(e){
    }
    Api.models().then((res)=>{
      const {code, data} = res
      if (code){
        return
      }
      const keys = {}
      const models = []
      const lables = t('setting.lables', { returnObjects: true })
      const alias = t('setting.alias', { returnObjects: true })
      Object.keys(data).map((key)=>{
        const node = data[key]
        if (!node.show){
          return
        }
        if(!keys[node.owner]){
            keys[node.owner] = key
            models.push({
                label: lables[node.owner],
                options:[{
                    label: alias[node.name] + " (" + (node.length / 1024) + "K)",
                    value: node.name,
                    disabled: node.disabled || !init[node.owner],
                    length: node.length
                }]
            })
            return
        }
        const item = models.find((m)=> {
            return lables[node.owner] == m.label
        })
        if(item){
            item.options.push({
                label: alias[node.name] + " (" + (node.length / 1024) + "K)",
                value: node.name,
                disabled: node.disabled || !init[node.owner],
                length: node.length
            })
        }
      })
      setModels(models)
    })
  }
  const closeTour = () => {
    localStorage.setItem('setting.tour', 1)
    setOpen(false)
  }

  const finish = async (vals) => {
    const items = {
      ...setting,
      ...vals
    }
    setDisable(true)
    const res = await Api.update_setting(bot.id, items)
    const { success, data } = res
    if (success) {
      setSetting(data)
      bot.setting = data
      message.success(t('save_success'))
    }
    setDisable(false)
  }

  const update = (key, val) => {
    setSetting({
      ...setting,
      [key]: val
    })
  }

  const scoreMarks = {};
  const tempMarks = {};
  const contexts = t('setting.contexts', { returnObjects: true }).map((val, idx)=>{return {label: val, value: idx}})
  const temperatures = t('setting.temperatures', { returnObjects: true })
  const model_tip = (<>{t('setting.model_tip', { returnObjects: true }).map((val, idx)=><div key={idx}>{val}</div>)}</>)
  const alert_tip = (<>{t('setting.alert_tip', { returnObjects: true }).map((val, idx)=><div key={idx}>{val}</div>)}</>)
  const formatter = (val) => {return temperatures[val]}
  ([50,60,70,80,90]).forEach((v) => {scoreMarks[v] = v});
  ([0,25,50,75,100]).forEach((v) => {tempMarks[v] = v});

  const rules = [{ required: true, message: t('required') }]

  return (
    <div className='settings-container px-5 pt-4'>
      <Contact
        show={visible}
        onChange={(show) => {
          setVisible(show)
        }}
        title={t('setting.modal.title')}
        description={t('setting.modal.desc')}/>
      <div>
        <Form form={form} {...layout} layout='vertical' name='setting' initialValues={setting} onFinish={finish} style={{ marginLeft: '18px' }}>
          <h5 className='pb-2'>{t('setting.profil')}</h5>
          <Form.Item style={{ marginBottom: 'auto' }}>
            <Row wrap={false}>
              <Col flex='80px'>
                <Avatar src={setting.avatar} size={64}></Avatar>
              </Col>
              <Col flex='auto'>
                <Form.Item name='name' style={{ marginBottom: '8px' }} rules={rules}>
                  <Input maxLength={20} />
                </Form.Item>
                <div className='sub-title'>ID: {bot.id}</div>
              </Col>
            </Row>
          </Form.Item>
          <h5 className='py-2'>{t('setting.description')}</h5>
          <Form.Item name='description' rules={rules}>
            <Input.TextArea placeholder='' showCount maxLength={120} />
          </Form.Item>
          <h5 className='py-2'>
            {t('setting.model')}
            {' '}
            <Tooltip placement='right' title={model_tip} overlayInnerStyle={{width: "29rem"}}>
              <QuestionCircleOutlined  ref={ref1}/>
            </Tooltip>
          </h5>
          <Form.Item name='model' rules={rules}>
            <Select
                defaultValue={setting.model}
                style={{ width: "100%" }}
                onChange={(val)=> console.log(val)}
                options={models}
              />
          </Form.Item>
          {
            app.language == "zh-CN" ? <Alert banner message={alert_tip} type="warning" /> : ""
          }
          <Form.Item style={{ marginBottom: 'auto' }}>
            <h5 className='py-2'>
              {t('setting.contexts_info')} {' '}
              <Tooltip placement='right' title={t('setting.contexts_title')}>
                <QuestionCircleOutlined />
              </Tooltip>
            </h5>
              <Form.Item name='contexts' rules={rules} >
                <Select
                  defaultValue={setting.contexts}
                  style={{ width: "100%" }}
                  onChange={(val)=> console.log(val)}
                  options={contexts}
                />
              </Form.Item>
          </Form.Item>
          <h5 className='py-2'>
            {t('setting.welcome')}{' '}
            <Tooltip placement='right' title={t('setting.welcome_title')}>
              <QuestionCircleOutlined />
            </Tooltip>{' '}
          </h5>
          <Form.Item name='welcome' rules={rules} extra={t('setting.prompt_extra')}>
            <Input.TextArea placeholder='' showCount maxLength={200} style={{ height: 160 }} />
          </Form.Item>
          <h5 className='py-2'>
            {t('setting.prompt')}{' '}
            <Tooltip placement='right' title={t('setting.prompt_title')}>
              <QuestionCircleOutlined  />
            </Tooltip>
          </h5>
          <div ref={ref3}>
            <Form.Item name='prompt' >
                <Input.TextArea placeholder='' showCount maxLength={600} style={{ height: 240 }} />
            </Form.Item>
          </div>
          {
            <>
            <h5 className='py-2'>
              {t('setting.score_info')}{' '}
              <Tooltip placement='right' title={t('setting.score_title')}>
                <QuestionCircleOutlined />
              </Tooltip>
            </h5>
            <Form.Item  rules={rules}>
              <div className="slider-wrapper">
                  <Slider marks={scoreMarks} onChange={(val) => update("score",val)} defaultValue={setting.score} min={50} max={90} step={1}/>
              </div>
            </Form.Item>
            <h5 className='py-2'>
              {t('setting.temperature_info')}{' '}
              <Tooltip placement='right' title={t('setting.temperature_title')}>
                <QuestionCircleOutlined />
              </Tooltip>
            </h5>
            <Form.Item  rules={rules}>
              <div className="slider-wrapper" ref={ref4}>
                <Slider marks={tempMarks} onChange={(val) => update("temperature",val)} defaultValue={setting.temperature} min={0} max={100} step={25} tooltip={{formatter}}/>
              </div>
            </Form.Item>
            </>
          }
          <Form.Item>
            <Button type='primary' size="large" disable={disable ? disable : undefined} htmlType='submit'>
              {t('save_change')}
            </Button>
          </Form.Item>
        </Form>
      </div>
      <Tour open={open} 
        mask={false} 
        type="primary" 
        onClose={closeTour} 
        steps={steps}
        indicatorsRender={(current, total) => (
          <span>
            {current + 1} / {total}
          </span>
        )}
        />
    </div>
  )
}

const mapStateToProps = (state) => {
  return {
    bot: state.bot,
    user: state.user,
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Setting)
