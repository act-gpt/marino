import React, { useState, useEffect } from 'react'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import { useTranslation } from 'react-i18next'
import { Tabs, Timeline, message } from 'antd'
import { CopyOutlined } from '@ant-design/icons'

const Mp = (props) => {
  const { token, id } = props
  const { t, i18n } = useTranslation()
  const [messageApi, contextHolder] = message.useMessage()
  const data = [
    {
      key: '1',
      label: '菜单',
      children: 'Content of Tab Pane 1'
    },
    {
      key: '2',
      label: '关键字回复',
      children: 'Content of Tab Pane 2'
    }
  ]
  return (
    <>
      <div>
        {contextHolder}
        <Tabs defaultActiveKey='1' items={data} />
      </div>
    </>
  )
}

export default Mp
