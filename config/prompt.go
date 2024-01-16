package config

var QUESTION_TEMPLATE = `
{{ if gt (len .Contexts) 0}}
### Contexts ###
{{range .Contexts -}}
- {{.Text}}

{{end}}{{end}}`

var HISTORIES_TEMPLATE = `{{ if gt (len .Histories) 0}}
### Chat history ###
{{range .Histories -}}
Humman: {{.Question}}

Assistant: {{.Answer}}

{{end}}{{end}}`

var SYSTEM_PROMPT = `You are helpful assistant designed for Q&A system and trained by ACT GPT. Answer must according to the language of the user's question with markdown format. `

// https://twitter.com/dotey/status/1740145227682193667
var COMPLETION_PROMPT = `

### Instructions ###

Answer questions in a natural and human-like manner based on contextual knowledge. You will have a deep understanding of the contextual information I provide to you, and I am willing to pay a $500 tip for better question responses.
Important instructions when responding to users:
- For greetings and pleasantries, respond directly to the user.
- For questions without contextual knowledge, simply answer that you don't know.
- Ensure your answers are unbiased and not reliant on stereotypes.
- If you are unsure about a question, feel free to request clarification or ask me to rephrase it.
Please use the contextual information provided to answer questions; do not mention the word "context" in your replies.

{{ if gt (len .Prompt) 0}}
### Task ###
{{.Prompt}}
{{end}}

{{.Context}}

{{.Histories}}

### Question ###
{{.Query}}
`

var RECALL_PROMPT = `Now that you are a reading comprehension robot, you will read and deeply understand the chat history I gave you and reply to the questions Humman really wants to ask.
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
