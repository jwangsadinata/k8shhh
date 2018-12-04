package k8shhh

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"text/template"

	"github.com/joho/godotenv"
)

// Secret is the type containing the name and the underlying data
type Secret struct {
	Name string
	Data map[string]string
}

// Encoder is a type for function that encodes the given Secret
type Encoder func(Secret) ([]byte, error)

// jsonTemplate is the template struct for json encoding
type jsonTemplate struct {
	ApiVersion string            `json:"apiVersion"`
	Data       map[string]string `json:"data"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Type       string            `json:"type"`
}

const (
	// yamlTemplate is the template string for yaml encoding
	yamlTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: {{.Name}}
type: Opaque
data:
{{range $key, $value := .Data}}  {{$key}}: {{encode $value}}
{{end}}`
)

// Encode encodes the input based on the given encoder
func Encode(input io.Reader, encoder Encoder, name string) ([]byte, error) {
	data, err := godotenv.Parse(input)
	if err != nil {
		return nil, err
	}

	return encoder(Secret{name, data})
}

// EncodeJson encodes the secret and output it to a json format
func EncodeJson(secret Secret) ([]byte, error) {
	tmpl := jsonTemplate{
		ApiVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		Data:       make(map[string]string),
		Metadata:   make(map[string]string),
	}

	tmpl.Metadata["name"] = secret.Name

	for k, v := range secret.Data {
		tmpl.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}

	return json.MarshalIndent(tmpl, "", "\t")
}

// EncodeYaml encodes the secret and output it to a yaml format
func EncodeYaml(secret Secret) ([]byte, error) {
	t, err := template.New("k8s-secrets").Funcs(template.FuncMap{
		"encode": func(input string) string {
			return base64.StdEncoding.EncodeToString([]byte(input))
		}}).Parse(yamlTemplate)
	if err != nil {
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, secret)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
