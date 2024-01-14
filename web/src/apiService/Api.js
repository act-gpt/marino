import { request } from './request'
import env from './ENV'

const apiPath = env.api + '/dashboard/'

// eslint-disable-next-line import/no-anonymous-default-export
export default {
  UPLOAD_URL: apiPath + 'knowledges/upload',
  /**
   * 当前用户
   *
   *
   * **/
  me() {
    return request(apiPath + 'me', 'GET')
  },

  check(type, body) {
    return request(apiPath + 'check', 'POST', {type}, body)
  },

  banner (lang) {
    return request('https://cdn.act-gpt.com/config/banner' + lang + ".json", 'GET')
  },

  config() {
    return request(apiPath + 'config', 'GET')
  },

  save_config(body) {
    return request(apiPath + 'config', 'POST', {}, body)
  },

  integrity(url) {
    return request(apiPath + 'js/integrity', 'GET', {url})
  },

  login(body) {
    return request(apiPath + 'users/login', 'POST', {}, body)
  },

  signout() {
    return request(apiPath + 'users/logout', 'GET')
  },

  signup(body) {
    return request(apiPath + 'users/register', 'POST', {}, body)
  },

  verification(email) {
    return request(apiPath + 'verification', 'POST', {}, { email })
  },

  resend(email) {
    return request(apiPath + 'reset_password', 'POST', {}, { email })
  },

  reset(data) {
    return request(apiPath + 'user/reset', 'POST', {}, data)
  },

  bots() {
    return request(apiPath + 'bots/', 'GET')
  },

  bot(id) {
    return request(apiPath + 'bots/' + id, 'GET')
  },

  query(id, data) {
    return request(apiPath + 'query/' + id, 'POST', {}, data)
  },

  datail(id, doc) {
    return request(apiPath + 's/' + id, 'GET', { id: doc })
  },

  conversation(id, data) {
    return request(apiPath + 'web/conversation/' + id, 'GET', data)
  },

  messages(id, data) {
    return request(apiPath + 'bots/messages/' + id, 'GET', data)
  },

  org() {
    return request(apiPath + 'orgs/', 'GET')
  },

  quato() {
    return request(apiPath + 'orgs/quato', 'GET')
  },

  models() {
    return request(apiPath + 'models', 'GET')
  },

  add_org(data) {
    return request(apiPath + 'orgs/', 'POST', {}, data)
  },

  add_bot(data) {
    return request(apiPath + 'bots/template', 'POST', {}, data)
  },

  update_bot(id, data) {
    return request(apiPath + 'bots/' + id, 'PUT', {}, data)
  },

  delete_bot(id) {
    return request(apiPath + 'bots/' + id, 'DELETE')
  },

  apps(lang) {
    return request(apiPath + 'templates', 'GET', { lang })
  },

  setting(id) {
    return request(apiPath + 'bots/setting/' + id, 'GET')
  },

  update_setting(id, data) {
    return request(apiPath + 'bots/setting/' + id, 'PUT', {}, data)
  },
  /**
   * 获取 folder
   * @param {*}
   * @returns
   */
  get_folders(bot_id) {
    return request(apiPath + 'folders/bot/' + bot_id, 'GET')
  },

  /**
   * 获取 folder
   * @param {*}
   * @returns
   */
  get_folder(id) {
    return request(apiPath + 'folders/' + id, 'GET')
  },

  /**
   * 添加 folder
   * @param {*} data
   * @returns
   */
  add_folder(data) {
    return request(apiPath + 'folders/', 'POST', {}, data)
  },

  // 删除 folder
  del_folder(id) {
    return request(apiPath + 'folders/' + id, 'DELETE', {})
  },

  // 更新 category
  put_folder(id, data) {
    return request(apiPath + 'folders/' + id, 'PUT', {}, data)
  },

  /**
   * 添加 knowledge
   * @param {*} data
   * @returns
   */
  add_knowledge(data) {
    return request(apiPath + 'knowledges/', 'POST', {}, data)
  },

  /**
   * 获取 knowledge
   * @param {*} token
   * @returns
   */
  get_knowledge(id, item) {
    return request(apiPath + 'knowledges/', 'GET', { ...item, bot: id })
  },

  // 更新 knowledge
  put_knowledge(id, data) {
    return request(apiPath + 'knowledges/' + id, 'PUT', {}, data)
  },

  // 删除 knowledge
  del_knowledge(id) {
    return request(apiPath + 'knowledges/' + id, 'DELETE')
  },

  // 批量删除 knowledge
  batch_del_knowledge(ids) {
    return request(apiPath + 'knowledges/', 'DELETE', {}, ids)
  }
}
