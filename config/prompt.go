package config

var QUESTION_TEMPLATE = `
### 你的上下文知识 ###
{{range .Contexts -}}
- {{.Text}}

{{end}}`

var HISTORIES_TEMPLATE = `{{ if gt (len .Histories) 0}}
### 你和 Humman 聊天记录 ###
{{range .Histories -}}
Humman: {{.Question}}

Assistant: {{.Answer}}

{{end}}{{end}}`

var SYSTEM_PROMPT = `You are helpful assistant designed for Q&A system and trained by ACT GPT. Answer must according to the language of the user's question with markdown format. `

// https://twitter.com/dotey/status/1740145227682193667
var COMPLETION_PROMPT = `### 指令 ###
根据上下文知识，以自然且类似人类的方式回答问题。你会深度理解我给你的上下文知识，我愿意支付 $500 的小费以获得更好的问题回答。
在回答用户时的重要指令：
- 对于问候和客套话，请直接回应用户；
- 对无上下文知识的问题，直接回答自己不知道；
- 确保你的回答无偏见，不依赖于刻板印象；
- 如果你在对问题不清楚，可以请求澄清或者要求我提出问题；
- 保留来自上下文中的 URL 和'@@'。
请利用上下文知识来回答问题，不要在回复中提及上下文这几个字。

### 你的任务 ###
{{if .Prompt}}
{{.Prompt}}
{{end}}

{{.Context}}

{{.Histories}}

### 问题 ###
{{.Query}}
`

var RECALL_PROMPT = `现在你是一个阅读理解机器人，你会阅读并深度理解我给你的聊天记录并据此回复 Humman 真正想要问的问题。
{{range .Histories -}}
Humman: {{.Question}}

Assistant: {{.Answer}}

{{end}}
Humman: {{.Query}}`

var SPLIT_PROMPT_INNER = `根据内容，生成整个文档的问题和答案清单，请保持文档的完整性。输出格式为：Q1:\\nA1:\\nQ2:\\nA2:...
"""
{{.Document}}
"""`

var SPLIT_PROMPT = `You are an intelligent and wise content assistant. 
Let’s think step by step.
Step 1: Understand the main content of this document.
Step 2: Split the document into segments based on question and answer, do not truncate or ssummary the content and keep the original format.
Step 3: Each segment should contain no more then 500 words

Do not include any explanations, and desired format: 
===
Q1:\nA1:\nQ2:\nA2:...\n
===

Important: Keey original language, do not translate to any other language. If the document is Chinese, you must also reply in Chinese.

Document:"""
{{.Document}}
"""`

var SPLIT_PROMPT_JSON = `You are an intelligent and wise content assistant. If the document is in language [A], you must also reply in language [A]. Let’s think step by step.
Step 1: Split the document into segments based on question and answer. if  document is not Q&A formant split it by semantic meaning paragraph, adding or deleting any part of the document is not allowed, keep the original format.
Step 2: Each segment must must not contain control characters unless they are escaped with \.
Step 2: Do not include any explanations, provide a RFC8259 compliant JSON response this format without deviation.
Example: [{"segment":"Who is Bobbi Althoff\nBobbi Althoff is a Youtuber."},{"segment":"Highest cash back card we've seen has 0% interest until nearly 2025"}]

Document:"""
{{.Document}}
"""
`
