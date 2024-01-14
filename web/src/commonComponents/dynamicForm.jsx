import React from 'react'
import { Form, Input, Checkbox, Select } from 'antd'
import PropType from 'prop-types'
import { MinusCircleOutlined } from '@ant-design/icons'
import validator from 'validator'
import { validatorRule } from '../utils/formValidator'

const { Option } = Select

const FormItem = Form.Item

const DynamicForm = ({
  index,
  initialValues,
  formName,
  nameKey,
  nameValue,
  keyValue,
  value,
  placeholderKey,
  placeholderValue,
  onValuesChange,
  labelKey,
  rulesKey,
  rulesValue,
  labelValue,
  deleteForm,
  selectValue,
  optionType,
  selectPlaceholder,
  rulesSelect,
  selectName,
  isCheckbox,
  onFinishFailed,
  form,
  selectOptionValue
}) => {
  // 校验数据结构attr_name
  const verify_dara_structure_attr_name = (rule, value) => {
    return new Promise(async (resolve, reject) => {
      if (value === undefined) {
        await reject(`数据结构不能为空`)
      } else if (validator.isEmpty(value)) {
        await reject(`数据结构不能为空`)
      } else if (validator.isLength(value, validatorRule.attr_name.minLen, validatorRule.attr_name.maxLen)) {
        await resolve()
      } else {
        await reject(`数据结构不能大于${validatorRule.attr_name.maxLen}个字符`)
      }
    })
  }
  return (
    <Form
      form={form}
      initialValues={initialValues}
      style={{ display: 'flex' }}
      name={formName}
      onValuesChange={(changedValues, allValues) => onValuesChange(index, allValues)}
      layout='inline'
      onFinishFailed={onFinishFailed}>
      <FormItem name={nameKey} label={labelKey} rules={[{ validator: verify_dara_structure_attr_name, required: true }]}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Input placeholder={placeholderKey} value={keyValue} />
        </div>
      </FormItem>
      <FormItem name={nameValue} label={labelValue} rules={rulesValue} style={{ marginBottom: 20 }}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Input placeholder={placeholderValue} value={value} />
        </div>
      </FormItem>
      <Form.Item name={selectName} label={selectValue} rules={rulesSelect}>
        {/* 要改变 Select 默认值 */}
        <Select style={{ width: 80 }} placeholder={selectPlaceholder} value={selectOptionValue}>
          {optionType.map((opt) => (
            <Option key={opt}>{opt}</Option>
          ))}
        </Select>
      </Form.Item>
      {isCheckbox ? null : (
        <Form.Item name='required' valuePropName='checked'>
          <Checkbox>是否必填</Checkbox>
        </Form.Item>
      )}
      <MinusCircleOutlined style={{ marginLeft: 3, marginTop: 9 }} onClick={deleteForm} />
    </Form>
  )
}

DynamicForm.propType = {
  initialValues: PropType.object,
  nameKey: PropType.string,
  nameValue: PropType.string,
  placeholderKey: PropType.string,
  rulesKey: PropType.array,
  rulesValue: PropType.array,
  labelValue: PropType.string,
  placeholderValue: PropType.string,
  formName: PropType.string,
  index: PropType.number,
  keyValue: PropType.string,
  value: PropType.string,
  onValuesChange: PropType.func,
  deleteForm: PropType.func,
  selectPlaceholder: PropType.string,
  selectName: PropType.string,
  optionType: PropType.array,
  selectValue: PropType.string,
  selectOptionValue: PropType.string
}

export default DynamicForm
