import React from 'react'
import { Modal, Button, Input } from 'antd'
import PropTypes from 'prop-types'

const DeleteModal = ({ isApplication, delete_cancel, delete_confirm, delete_title, visible, width, style, delete_app_name, delete_application, delete_common_name, app_name }) => {
  return (
    <Modal
      visible={visible}
      title={delete_title}
      onCancel={delete_cancel}
      onOk={delete_confirm}
      width={width}
      style={style}
      footer={[
        <Button onClick={delete_cancel}>取消</Button>,
        <Button type='danger' onClick={delete_confirm}>
          删除
        </Button>
      ]}>
      {isApplication ? (
        <React.Fragment>
          <span className='application-modal-title'>
            请输入应用名<span style={{ color: '#f50' }}>{delete_app_name}</span>进行确认
          </span>
          <div className='application-modal-item'>
            <span className='title'>应用名称:</span>
            <Input className='input' placeholder='应用名称' onChange={delete_application} allowClear={true} value={app_name} />
          </div>
        </React.Fragment>
      ) : (
        <>
          <span className='application-modal-title'>
            是否删除 <span style={{ color: '#f50' }}>{delete_common_name}</span>?
          </span>
        </>
      )}
    </Modal>
  )
}

DeleteModal.propTypes = {
  visible: PropTypes.bool,
  title: PropTypes.string
}

export default DeleteModal
