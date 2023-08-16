package template

import (
	"bytes"
	"testing"
	"text/template"
)

type MyObj struct {
	IamBaseUrl         string
	InternalIamBaseUrl string
}

func TestStringTemplate(t *testing.T) {
	var myObj = MyObj{
		"http://localhost:8080",
		"http://192.168.159.129:8080/",
	}

	tmpl, err := template.New("urlTemplate").Funcs(template.FuncMap{
		"value": func(baseUrl string, defaultUrl string) string {
			if len(baseUrl) > 0 {
				return baseUrl
			}
			return defaultUrl
		},
	}).Parse("{{ value .InternalIamBaseUrl .IamBaseUrl }}/realms/master")
	if err != nil {
		panic(err)
	}
	var tmpBytes bytes.Buffer
	err = tmpl.Execute(&tmpBytes, myObj)
	if err != nil {
		panic(err)
	}

	println(tmpBytes.String())
}
