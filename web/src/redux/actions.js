import { SET_USER, SET_BOT, SET_APP, SET_CONFIG, MENU_TOGGLE } from './actionTypes'

const setUser = (user) => ({ type: SET_USER, payload: { user } })
const setBot = (bot) => ({ type: SET_BOT, payload: { bot } })
const setApp = (app) => ({ type: SET_APP, payload: { app } })
const setConfig = (config) => ({ type: SET_CONFIG, payload: { config } })
const tiggerToggle = (toggle) => ({ type: MENU_TOGGLE, payload: { toggle } })

const model =  {
  setUser,
  setApp,
  setBot,
  setConfig,
  tiggerToggle
}

export default model