import React from 'react'
import { Modal, Button } from 'antd'

const ExtendModal = ({ extend_title, visible, extend_cancle, extend_confirm, width, children, style }) => {
  return (
    <Modal
      visible={visible}
      title={extend_title}
      onCancel={extend_cancle}
      onOk={extend_confirm}
      width={width}
      style={style}
      footer={[
        <Button onClick={extend_cancle}>取消</Button>,
        <Button type='primary' onClick={extend_confirm}>
          确定
        </Button>
      ]}>
      {children}
    </Modal>
  )
}

export default ExtendModal
