import _ from 'lodash'
import validator from 'validator'

export const validatorRule = {
  app_name: {
    name: '应用名称',
    minLen: 1,
    maxLen: 8
  },
  app_id: {
    name: '应用代号',
    minLen: 1,
    maxLen: 15
  },
  category: {
    name: '类型名称',
    minLen: 1,
    maxLen: 6
  },
  category_style: {
    name: '类型样式'
  },
  user_group_name: {
    name: '用户组名',
    minLen: 1,
    maxLen: 8
  },
  user_name: {
    name: '用户名',
    minLen: 1,
    maxLen: 8
  },
  user_password: {
    name: '密码',
    minLen: 6,
    maxLen: 15
  },
  user_auth: {
    name: '用户权限'
  },
  blockchain_symbol: {
    name: '区块链服务名称',
    minLen: 1,
    maxLen: 8
  },
  blockchain_url: {
    name: '区块链服务地址'
  },
  template_type: {
    name: '模板类别'
  },
  template_id: {
    name: '模板代号',
    minLen: 1,
    maxLen: 15
  },
  template_name: {
    name: '模板名称',
    minLen: 1,
    maxLen: 8
  },
  app_data_structure_name: {
    name: '应用数据结构名称',
    minLen: 1,
    maxLen: 12
  },
  app_data_structure_id: {
    name: '应用数据结构代号',
    minLen: 1,
    maxLen: 12
  },
  app_data_structure_edit_id: {
    name: '应用数据结构编辑代号',
    minLen: 1,
    maxLen: 30
  },
  business_flow_name: {
    name: '业务流程名称',
    minLen: 1,
    maxLen: 12
  },
  business_unit_name: {
    name: '业务单元名称',
    minLen: 1,
    maxLen: 15
  },
  business_unit_type: {
    name: '业务单元类型'
  },
  business_unit_temp: {
    name: '业务单元模板'
  },
  receiver_name: {
    name: '接收方'
  },
  relational_business: {
    name: '关联业务'
  },
  filter_name: {
    name: '过滤条件'
  },
  dchain_data_type: {
    name: '下链数据类型'
  },
  attr_name: {
    name: '数据结构/属性名',
    minLen: 1,
    maxLen: 8
  },
  attr_key: {
    name: '数据结构/键'
  },
  attr_type: {
    name: '数据结构/数据类型'
  },
  uchain_data_type: {
    name: '上链数据模版'
  },
  data_structure: {
    name: '数据结构'
  },
  relational_business_flow: {
    name: '业务流程'
  },
  relational_business_unit: {
    name: '业务单元'
  },
  data_sources: {
    name: '数据来源'
  },
  rule_code: {
    name: '规则代号',
    minLen: 1,
    maxLen: 20
  }
}

/**
 * 用于处理输入框长度验证
 * @param {*} value input框输入值
 * @param {*} item 验证项(validatorRule中选择)
 * @param {*} setState react的setState方法,改变当前input框验证状态(0失败，1成功)
 * @param {*} state 当前form验证状态的state object
 * @param {*} key 当前验证项在state中的的key
 */
export const input_len_validator = (value, item) => {
  return new Promise(async (resolve, reject) => {
    if (!value) {
      await reject(`${item.name}不能为空`)
    } else if (validator.isEmpty(value)) {
      await reject(`${item.name}不能为空`)
    } else if (validator.isLength(value, item.minLen, item.maxLen)) {
      await resolve()
    } else {
      if (value.length < item.minLen) {
        await reject(`${item.name}不能小于${item.minLen}个字符`)
      } else {
        await reject(`${item.name}不能大于${item.maxLen}个字符`)
      }
    }
  })
}

/**
 * 用于判断select下拉不能为空
 * @param {*} value
 * @param {*} item
 * @param {*} setState
 * @param {*} state
 * @param {*} key
 */
export const select_empty_validator = (value, item) => {
  //对"",[]等false值做判断
  if (value && !value.toString()) {
    value = undefined
  }
  return new Promise(async (resolve, reject) => {
    //对number等值做转换
    if (value !== undefined) {
      value = value.toString()
    }
    if (!value || validator.isEmpty(value)) {
      await reject(`${item.name}不能为空`)
    } else {
      await resolve()
    }
  })
}

/**
 * 用于校验添加规则-规则代号字段
 * @param {*} value
 * @param {*} item
 * @returns
 */
export const input_len_alphanumeric_validator = (value, item) => {
  return new Promise(async (resolve, reject) => {
    if (!value) {
      await reject(`${item.name}不能为空`)
    } else if (validator.isEmpty(value)) {
      await reject(`${item.name}不能为空`)
    } else if (!validator.isAlphanumeric(value)) {
      await reject(`${item.name}不能为中文`)
    } else if (validator.isLength(value, item.minLen, item.maxLen)) {
      await resolve()
    } else {
      if (value.length < item.minLen) {
        await reject(`${item.name}不能小于${item.minLen}个字符`)
      } else {
        await reject(`${item.name}不能大于${item.maxLen}个字符`)
      }
    }
  })
}

/**
 * 校验区块链地址
 * @param {*} value
 * @param {*} item
 * @returns
 */
export const url_format_validator = (value, item) => {
  return new Promise(async (resolve, reject) => {
    if (!value || validator.isEmpty(value)) {
      await reject(`${item.name}不能为空`)
    } else if (validator.isURL(value)) {
      await resolve()
    } else {
      await reject(`${item.name}格式错误`)
    }
  })
}

/**
 * @param {Object} args
 * @returns 成功 1, 失败 0
 */
export const isValidForm = (args) => {
  return _.reduce(
    args,
    (pre, cur) => {
      return pre && cur
    },
    1
  )
}
