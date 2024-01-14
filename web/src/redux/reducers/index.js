import { combineReducers } from "redux";
import user from "./user";
import bot from "./bot"
import app from "./app"
import toggle from "./menuToggle"
import config from "./config"
export default combineReducers({user, bot, app, toggle, config});