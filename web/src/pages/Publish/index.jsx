import React, { useState } from 'react'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import { Col, Row, Card, Space, Button, message } from 'antd'
import { useTranslation } from 'react-i18next'
import { ContainerOutlined } from '@ant-design/icons'
import {
  GlobalOutlined,
  CodeSandboxOutlined,
  CopyOutlined,
  EyeOutlined,
  QqOutlined,
  WhatsAppOutlined,
  CommentOutlined,
  FacebookOutlined,
  DeploymentUnitOutlined,
  createFromIconfontCN
} from '@ant-design/icons'

import CommonModal from '@/commonComponents/commonModal'

import { connect } from 'react-redux'
import actions from '@/redux/actions'
import Javascript from './javascript'
import Api from './api'
const IconFont = createFromIconfontCN({
  scriptUrl: '//at.alicdn.com/t/c/font_4213103_cfasccvi9zl.js'
})

const Publish = (props) => {
  const { bot, app } = props
  const { t } = useTranslation()
  const [messageApi, contextHolder] = message.useMessage()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [model, setModel] = useState(1)

  const showModal = () => {
    setIsModalOpen(true)
  }

  const handleCancel = () => {
    setIsModalOpen(false)
  }
  const webPath = (path) => {
    return (process.env.REACT_APP_WEB || window.location.origin) + path
  }

  const viewPage = () => {
    window.open(webPath('/chat/' + bot.id), '_blank').focus()
  }

  return (
    <div className='settings-container px-5 pt-4'>
      <CommonModal
        title={model === 1 ? t('publish.js.modal.title') : model === 2 ? t('publish.api.modal.title') : ''}
        width={600}
        visible={isModalOpen}
        footer={false}
        common_cancel={handleCancel}
        children={model === 1 ? <Javascript host={process.env.REACT_APP_CDN || window.location.origin} id={bot.id} /> : model === 2 ? <Api token={bot.access_token} id={bot.id} /> : ''}
      />
      {contextHolder}
      <div>
        <Row gutter={[16, 16]}>
          <Col span={12}>
            <Card>
              <h4 className='pb-3'>
                <GlobalOutlined /> {t('publish.page.title')}
              </h4>
              <div className='pb-3 sub-title'> {t('publish.page.subtitle')}</div>
              <div className='pb-3'>
                <Space>
                  <CopyToClipboard
                    text={webPath('/chat/' + bot.id)}
                    onCopy={() =>
                      messageApi.open({
                        type: 'success',
                        content: t('copy_success')
                      })
                    }
                  >
                    <Button type='text' icon={<CopyOutlined />}>
                      {' '}
                      {t('publish.page.copy')}
                    </Button>
                  </CopyToClipboard>
                  <Button type='text' onClick={viewPage} icon={<EyeOutlined />}>
                    {t('publish.page.view')}
                  </Button>
                </Space>
              </div>
            </Card>
          </Col>
          <Col span={12}>
            <Card>
              <h4 className='pb-3'>
                <CodeSandboxOutlined /> {t('publish.js.title')}
              </h4>
              <div className='pb-3 sub-title'> {t('publish.js.subtitle')} </div>
              <div className='pb-3'>
                <Space>
                  <Button
                    type='text'
                    onClick={() => {
                      setModel(1)
                      showModal()
                    }}
                    icon={<EyeOutlined />}
                  >
                    {t('publish.js.view')}
                  </Button>
                </Space>
              </div>
            </Card>
          </Col>
        </Row>
        <Row className='mt-5' gutter={[16, 16]}>
          <Col span={12}>
            <Card>
              <h4 className='pb-3'>
                <DeploymentUnitOutlined /> {t('publish.api.title')}
              </h4>
              <div className='pb-3 sub-title'> {t('publish.api.subtitle')}</div>
              <div className='pb-3'>
                <Space>
                  <Button
                    type='text'
                    onClick={() => {
                      setModel(2)
                      showModal()
                    }}
                    icon={<EyeOutlined />}
                  >
                    {t('publish.api.view')}
                  </Button>
                  <Button
                    type='text'
                    onClick={() => {
                      window.open(t('publish.api.url'), '_blank').focus()
                    }}
                    icon={<ContainerOutlined />}
                  >
                    {' '}
                    {t('publish.api.document')}
                  </Button>
                </Space>
              </div>
            </Card>
          </Col>
          {/* 
            <Col span={12}>
            <Card>
              <h4 className='pb-3'>
                <WechatOutlined /> 公众号
              </h4>
              <div className='pb-3 sub-title'> {t('publish.api.subtitle')}</div>
              <div className='pb-3'>
                <Space>
                  <Button
                    type="text"
                    onClick={() => {
                      setModel(2)
                      showModal()
                    }}
                    icon={<EyeOutlined />}>
                    {t('publish.api.view')}
                  </Button>
                  <Button type="text" onClick={() => {
                    window.open( t("publish.api.url"), '_blank').focus()
                  }} icon={<ContainerOutlined />}> {t('publish.api.document')}</Button>
                </Space>
              </div>
            </Card>
          </Col>
          */}
        </Row>
      </div>
      <div className='text-center mt-5'>
        <div>
          <h5 className='pb-3'>{t('publish.footer.comingsoon')}</h5>
          <Space>
            {app.language === 'en' ? (
              <>
                <div>
                  <div>
                    <WhatsAppOutlined />
                  </div>
                  <div className='sub-title'>{t('publish.footer.whatsapp')}</div>
                </div>
                <div>
                  <div>
                    <FacebookOutlined />
                  </div>
                  <div className='sub-title'>{t('publish.footer.facebook')}</div>
                </div>
              </>
            ) : (
              <>
                <div>
                  <div>
                    <CommentOutlined />
                  </div>
                  <div className='sub-title'>{t('publish.footer.wechat_group')}</div>
                </div>
                <div>
                  <div>
                    <IconFont type={'icon-xiaochengxu'} />
                  </div>
                  <div className='sub-title'>{t('publish.footer.wechat_program')}</div>
                </div>
                <div>
                  <div>
                    <IconFont type={'icon-contract'} />
                  </div>
                  <div className='sub-title'>{t('publish.footer.wechat_custome')}</div>
                </div>
                <div>
                  <div>
                    <QqOutlined />
                  </div>
                  <div className='sub-title'>{t('publish.footer.qq')}</div>
                </div>
              </>
            )}
          </Space>
        </div>
      </div>
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

export default connect(mapStateToProps, mapDispatchToProps)(Publish)
