package testing

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/lamoda/gonkey/models"
)

type TestingOutput struct {
	testing *testing.T
}

func NewOutput(t *testing.T) *TestingOutput {
	return &TestingOutput{
		testing: t,
	}
}

func (o *TestingOutput) Process(t models.TestInterface, result *models.Result) error {
	if !result.Passed() {
		text, err := renderResult(result)
		if err != nil {
			return err
		}
		o.testing.Error(text)
	}
	return nil
}

func renderResult(result *models.Result) (string, error) {
	text := `
       Name: {{ .Test.GetName }}
       File: {{ .Test.GetFileName }}

Request:
     Method: {{ .Test.GetMethod }}
       Path: {{ .Test.Path }}
      Query: {{ .Test.ToQuery }}
{{- if .Test.Headers }}
    Headers: 
{{- range $key, $value := .Test.Headers }}
      {{ $key }}: {{ $value }}
{{- end }}
{{- end }}
{{- if .Test.Cookies }}
    Cookies: 
{{- range $key, $value := .Test.Cookies }}
      {{ $key }}: {{ $value }}
{{- end }}
{{- end }}
       Body:
{{ if .RequestBody }}{{ .RequestBody }}{{ else }}{{ "<no body>" }}{{ end }}

Response:
     Status: {{ .ResponseStatus }}
       Body:
{{ if .ResponseBody }}{{ .ResponseBody }}{{ else }}{{ "<no body>" }}{{ end }}

{{ range $i, $dbr := .DatabaseResult }}
{{ if $dbr.Query }}
       Db Request #{{ inc $i }}:
{{ $dbr.Query }}
       Db Response #{{ inc $i }}:
{{ range $value := $dbr.Response }}
{{ $value }}{{ end }}
{{ end }}
{{ end }}

{{ if .Errors }}
     Result: {{ "ERRORS!" }}

Errors:
{{ range $i, $e := .Errors }}
{{ inc $i }}) {{ $e.Error }}
{{ end }}
{{ else }}
     Result: {{ "OK" }}
{{ end }}
`

	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}

	var buffer bytes.Buffer
	t := template.Must(template.New("letter").Funcs(funcMap).Parse(text))
	if err := t.Execute(&buffer, result); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
