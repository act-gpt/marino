import { SET_BOT } from "../actionTypes";

const initialState = {
};

const model = (state = initialState, action) => {
  switch (action.type) {
    case SET_BOT: {
      const { bot } = action.payload;
      return {
        ...state,
        ...bot,
      };
    }
    default:
      return state;
  }
}
export default model