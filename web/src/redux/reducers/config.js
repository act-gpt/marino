import { SET_CONFIG } from "../actionTypes";

const initialState = {
    Organization: {
      Name: "Marino"
    },
    Embedding: {
      Host: "https://api.openai.com",
      Api: "/v1/embeddings",
      Model: "text-embedding-ada-002",
    },
    Db: {
      DataSource: "",
      Dimension: 768
    },
    Moderation: {
      Api: "",
      CheckContent: false,
    },
    ActGpt: {
      AccessKey: ""
    },
    Baidu: {
        ClientId: "",
        ClientSecret: ""
    },
    OpenAi: {
      Type: "openai",
      Host: "https://api.openai.com",
      APIVersion: "2023-05-15",
  }
}
const model = (state = initialState, action) => {
  switch (action.type) {
    case SET_CONFIG: {
      return {
        ...state,
        ...action.payload.config,
      };
    }
    default:
      return state;
  }
}
export default model