package common

import (
	"bytes"
	"text/template"
)

type PromptTemplate string

func (p PromptTemplate) Render(data any) (string, error) {
	tmpl, err := template.New("").Parse(string(p))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
