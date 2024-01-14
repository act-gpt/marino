import { SET_USER } from "../actionTypes";

const initialState = {
};

const model = (state = initialState, action) => {
  switch (action.type) {
    case SET_USER: {
      const { user } = action.payload;
      return {
        ...state,
        ...user,
      };
    }
    default:
      return state;
  }
}

export default model