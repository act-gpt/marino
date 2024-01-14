import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import actions from '@/redux/actions'
import CommonModal from '@/commonComponents/commonModal'

const Contact = (props) => {
  const { title, description, onChange } = props
  const { t, i18n } = useTranslation()
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    setVisible(props.show)
  }, [props.show])

  return (
    <CommonModal
      title={title}
      width={600}
      visible={visible}
      common_cancel={() => {
        onChange && onChange(!props.show)
        setVisible(false)
      }}
      footer={false}
      children={
        props.app.language === 'en' ? (
          <div className='text-center'>
            <p>
              <p>{t('setting.modal.comingsoon')}</p>
            </p>
          </div>
        ) : (
          <div className='text-center'>
            {description ? <p>{description}</p> : ''}
            {
              <p className='mt-3'>
                <img src='/imgs/contactor.jpg' width='180' />
              </p>
            }
          </div>
        )
      }
    />
  )
}
const mapStateToProps = (state) => {
  return {
    user: state.user,
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}
export default connect(mapStateToProps, mapDispatchToProps)(Contact)
