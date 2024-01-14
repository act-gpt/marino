import { SET_APP } from "../actionTypes";

const initialState = {
  loading: true,
  auth: true,
};

const model = (state = initialState, action) => {
  switch (action.type) {
    case SET_APP: {
      return {
        ...state,
        ...action.payload.app,
      };
    }
    default:
      return state;
  }
}
export default model