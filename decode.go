package k8shhh

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// Decoder is a type for function that decodes a given io.Reader input
type Decoder func(io.Reader) (interface{}, error)

// Decode decodes the input based on the given decoder
func Decode(input io.Reader, decoder Decoder) ([]byte, error) {
	res, err := decoder(input)
	if err != nil {
		return []byte{}, err
	}

	var secret map[string]interface{}

	switch res := res.(type) {
	case map[interface{}]interface{}:
		secret = convertKeysToStrings(res)
	case map[string]interface{}:
		secret = res
	default:
		return []byte{}, errors.New(fmt.Sprintf("unexpected type: %T", res))
	}

	d := secret["data"]
	if d == nil {
		return []byte{}, nil
	}

	var data map[string]interface{}

	switch d := d.(type) {
	case map[interface{}]interface{}:
		data = convertKeysToStrings(d)
	case map[string]interface{}:
		data = d
	default:
		return []byte{}, errors.New(fmt.Sprintf("unexpected type: %T", d))
	}

	processed := convertValuesToStrings(data)

	for k, v := range processed {
		l, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return []byte{}, err
		}
		processed[k] = string(l)
	}

	lines := make([]string, 0, len(processed))
	for k, v := range processed {
		lines = append(lines, fmt.Sprintf(`%s=%s`, k, doubleQuoteEscape(v)))
	}
	sort.Strings(lines)
	out := strings.Join(lines, "\n")

	return []byte(out), nil
}

// DecodeJson decodes the json formatted input into the readable secret
func DecodeJson(input io.Reader) (interface{}, error) {
	var res interface{}
	if err := json.NewDecoder(input).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// DecodeYaml decodes the yaml formatted input into the readable secret
func DecodeYaml(input io.Reader) (interface{}, error) {
	var res interface{}
	if err := yaml.NewDecoder(input).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}

// convertKeysToStrings converts the keys of a given map to strings
func convertKeysToStrings(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	for k, v := range m {
		res[fmt.Sprintf("%v", k)] = v
	}

	return res
}

// convertValuesToStrings converts the values of a given map to strings
func convertValuesToStrings(m map[string]interface{}) map[string]string {
	res := make(map[string]string)

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

// doubleQuoteEscape is a helper function for escaping double quotes
func doubleQuoteEscape(line string) string {
	for _, c := range "\\\n\r\"!$`" {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}
