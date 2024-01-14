import React from 'react'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import { useTranslation } from 'react-i18next'
import { List, message } from 'antd'
import { CopyOutlined } from '@ant-design/icons'

const Api = (props) => {
  const { token, id } = props
  const { t } = useTranslation()
  const [messageApi, contextHolder] = message.useMessage()
  const data = [
    {
      name: t('publish.api.modal.id'),
      value: (process.env.REACT_APP_WEB || window.location.origin ) + '/v1/chat/'+ id,
      desc: t('publish.api.modal.bot_id')
    },
    {
      name: t('publish.api.modal.token'),
      value: token,
      desc: t('publish.api.modal.bot_token')
    }
  ]
  return (
    <>
      <div>
        {contextHolder}
        <List
          itemLayout='horizontal'
          dataSource={data}
          renderItem={(item, index) => (
            <List.Item>
              <List.Item.Meta title={item.name} description={item.desc} />
              <CopyToClipboard
                text={item.value}
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
            </List.Item>
          )}
        />
      </div>
    </>
  )
}

export default Api
