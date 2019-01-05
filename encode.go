package k8shhh

import (
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// Secret is the type containing the name and the underlying data
type Secret struct {
	Name string
	Data map[string]string
}

// Encoder is a type for function that encodes the given Secret
type Encoder func(Secret) ([]byte, error)

// template is the template struct for both json and yaml encoding
type template struct {
	APIVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Data       map[string]string `json:"data" yaml:"data"`
	Kind       string            `json:"kind" yaml:"kind"`
	Metadata   map[string]string `json:"metadata" yaml:"metadata"`
	Type       string            `json:"type" yaml:"type"`
}

// Encode encodes the input based on the given encoder
func Encode(input io.Reader, encoder Encoder, name string) ([]byte, error) {
	data, err := godotenv.Parse(input)
	if err != nil {
		return nil, err
	}
	return encoder(Secret{name, data})
}

// EncodeJSON encodes the secret and output it to a json format
func EncodeJSON(secret Secret) ([]byte, error) {
	return json.MarshalIndent(generateTemplate(secret), "", "\t")
}

// EncodeYAML encodes the secret and output it to a yaml format
func EncodeYAML(secret Secret) ([]byte, error) {
	return yaml.Marshal(generateTemplate(secret))
}

// generateTemplate puts the Secret name and data to the kubernetes template
func generateTemplate(secret Secret) template {
	tmpl := template{
		APIVersion: "v1",
		Data:       make(map[string]string),
		Kind:       "Secret",
		Metadata:   map[string]string{"name": secret.Name},
		Type:       "Opaque",
	}
	for k, v := range secret.Data {
		tmpl.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}
	return tmpl
}
