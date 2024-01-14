import React, { useState, useEffect } from 'react'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import { useTranslation } from 'react-i18next'
import { Card, message } from 'antd'
import Api from '@/apiService/Api'
import { CopyOutlined } from '@ant-design/icons'

const Javascript = (props) => {
  const { host, id } = props
  const { t } = useTranslation()
  const [src, setSrc] = useState('')
  const [messageApi, contextHolder] = message.useMessage()
  const update = async () => {
    const url = `${host}/js/embed.js`
      setSrc(`(function(){if(!window.__act_gpt){window.__act_gpt={"id": "${id}","user": ""}}var ele = document.createElement('script');ele.src = '${url}';document.body.appendChild(ele)}())`)
  }
  useEffect(() => {
    update()
  }, [id])

  return (
    <>
      <div>
        {contextHolder}
        <Card
          title={t('publish.js.modal.card_title')}
          extra={
            <CopyToClipboard
              text={src}
              onCopy={() =>
                messageApi.open({
                  type: 'success',
                  content: t('copy_success')
                })
              }>
                <span className='anticon' tabIndex={-1}>
                  <CopyOutlined />
                  <span className="ml-1">{t('publish.js.modal.copy')}</span>
                </span>
            </CopyToClipboard>
          }>
          <code>{src}</code>
        </Card>
      </div>
      <div className='pt-2 sub-title' dangerouslySetInnerHTML={{ __html: t('publish.js.modal.tip') }}></div>
    </>
  )
}

export default Javascript
