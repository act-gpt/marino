
import { MENU_TOGGLE } from "../actionTypes";

const initialState = {
  menuToggle: false
};

const model = (state = initialState, action) => {
  switch (action.type) {
    case MENU_TOGGLE:
      return { ...state, menuToggle: !state.menuToggle }
    default:
      return state
  }
}

export default model