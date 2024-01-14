import React from 'react'
import { Form, Input } from 'antd'
import PropType from 'prop-types'

const VerificationForm = ({ verifyValidator, name, form, onChange, placeholder, allowClear, label, style, type, layout, initialValue, maxLength }) => {
  return (
    <Form form={form} layout='horizontal' initialValues={initialValue} onValuesChange={onChange} {...layout}>
      <Form.Item name={name} style={style} label={label} rules={[{ validator: verifyValidator }]}>
        <Input type={type} style={{ 'outline-color': 'red' }} placeholder={placeholder} allowClear={allowClear} maxLength={maxLength} />
      </Form.Item>
    </Form>
  )
}

VerificationForm.propType = {
  verifyValidator: PropType.func.isRequired, // 校验方法
  name: PropType.string, // form item name
  onChange: PropType.func.isRequired, // input change 回调方法
  value: PropType.string, // input value
  placeholder: PropType.string,
  allowClear: PropType.bool,
  style: PropType.object,
  type: PropType.string,
  label: PropType.string
}

export default VerificationForm
