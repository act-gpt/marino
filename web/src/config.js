const config = {
  bot: {
    settings: {
      prompt: `你是一位管理公司行政事务的秘书，负责描述和解释公司相关规定。请根据下面内容回答用户问题，让用户能清晰明白公司的各种制度和流程。如果提供的信息没有和问题相关，回答："没有找到相关制度，如果您还有其他问题或需要进一步解答，请随时提出。"。用户的问题由 #### 确定

        {{range .Sections -}}
        * {{.}}
        {{- end}}
        
        ####
        {{.Question}}
        ####
        `
    },
    config: {
      "collectionName": "doc",
      "corpus": "actgpt:test",
      "dimension": 768,
      "indexType": "IVF_FLAT",
      "metricType": "IP",
      "model": "m3e",
      "shards": 0,
      "topK": 3,
      "units": 3
    }
  },
  org: {
    infomation: {},
    config: {
      "limited": 1
    },
  }
}

module.exports = config
