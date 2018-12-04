package k8shhh

import (
	"bytes"
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

	var data map[string]interface{}

	switch res := res.(type) {
	case map[interface{}]interface{}:
		data = convertKeysToStrings(res)
	case map[string]interface{}:
		data = res
	default:
		return []byte{}, errors.New(fmt.Sprintf("unexpected type: %T", res))
	}

	n := data["data"]
	var d map[string]interface{}

	switch n := n.(type) {
	case map[interface{}]interface{}:
		d = convertKeysToStrings(n)
	case map[string]interface{}:
		d = n
	default:
		return []byte{}, errors.New(fmt.Sprintf("unexpected type: %T", n))
	}

	fin := convertValuesToStrings(d)

	for k, v := range fin {
		l, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return []byte{}, err
		}
		fin[k] = string(l)
	}

	lines := make([]string, 0, len(fin))
	for k, v := range fin {
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

	//	b, err := readFile(input)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if err := yaml.Unmarshal(b, &res); err != nil {
	//		return nil, err
	//	}
	return res, nil
}

// readFile is a helper function for converting io.Reader to []byte
func readFile(input io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, input)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func convertKeysToStrings(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	for k, v := range m {
		res[fmt.Sprintf("%v", k)] = v
	}

	return res
}

func convertValuesToStrings(m map[string]interface{}) map[string]string {
	res := make(map[string]string)

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

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
