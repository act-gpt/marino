import React from 'react'
import { Modal, Button } from 'antd'
import '@/style/modal.scss'

import { useTranslation} from 'react-i18next'

const CommonModal = (props) => {
  const { t, i18n } = useTranslation();
  const { title, visible, centered, common_cancel, common_confirm, width, children, danger, footer=true, custom_foot=""} = props

  return (
    <Modal forceRender={true} getContainer={false} closable={false} width={width} visible={visible} centered={centered} footer={null}>
      <div className='modal-container'>
        <div className='modal-header'>
          <span />
          <span className='title'>{title}</span>
          <span onClick={common_cancel} className='icon'>
            X
          </span>
        </div>
        <div className='modal-content'>{children}</div>
        {
          footer ?
            <div className='modal-footer'>
              <Button className='btn' type='default' onClick={common_cancel}>
                {t('cancel')}
              </Button>
              <span className='space' />
              <Button className='btn' type='primary' danger={danger} onClick={common_confirm}>
              {t('ok')}
              </Button>
          </div>
          : ""
        }
        {
          custom_foot
        }
      </div>
    </Modal>
  )
}

export default CommonModal
