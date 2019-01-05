package k8shhh

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

// TestDecode tests the Decode function
func TestDecode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   io.Reader
		decoder Decoder
		name    string
		res     string
		err     error
	}{
		{
			input:   strings.NewReader("-1"),
			decoder: DecodeYAML,
			name:    "error-test",
			err:     errors.New("unexpected type: int"),
		},
		{
			input:   strings.NewReader("data: -1"),
			decoder: DecodeYAML,
			name:    "error-test2",
			err:     errors.New("unexpected type: int"),
		},
		{
			input:   strings.NewReader("value: -"),
			decoder: DecodeYAML,
			name:    "error-test3",
			err:     errors.New("yaml: block sequence entries are not allowed in this context"),
		},
		{
			input:   strings.NewReader(errorDecodeYAMLTest),
			decoder: DecodeYAML,
			name:    "error-test4",
			err:     errors.New("illegal base64 data at input byte 0"),
		},
		{
			input:   strings.NewReader(successDecodeJSONTestEmpty),
			decoder: DecodeJSON,
			name:    "json-empty",
			res:     "",
		},
		{
			input:   strings.NewReader(successDecodeYAMLTestEmpty),
			decoder: DecodeYAML,
			name:    "yaml-empty",
			res:     "",
		},
		{
			input:   strings.NewReader(successDecodeYAMLTestEmpty2),
			decoder: DecodeYAML,
			name:    "yaml-empty2",
			res:     "",
		},
		{
			input:   strings.NewReader(successDecodeJSONTestOne),
			decoder: DecodeJSON,
			name:    "json-one",
			res:     "a=b",
		},
		{
			input:   strings.NewReader(successDecodeYAMLTestOne),
			decoder: DecodeYAML,
			name:    "yaml-one",
			res:     "a=b",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := Decode(test.input, test.decoder)
			if err == nil {
				if test.err != nil {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			} else {
				if err.Error() != test.err.Error() {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			}
			if string(res) != test.res {
				t.Fatalf("expected response to be %q but got %q", test.res, res)
			}
		})
	}
}

// TestDecodeJSON tests the DecodeJSON function
func TestDecodeJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input io.Reader
		name  string
		res   interface{}
		err   error
	}{
		{
			input: strings.NewReader(`[asdf]`),
			name:  "decode-json-error",
			err:   errors.New("invalid character 'a' looking for beginning of value"),
		},
		{
			input: strings.NewReader(`{"a":"b"}`),
			name:  "decode-json-success",
			res:   map[string]interface{}{"a": "b"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := DecodeJSON(test.input)
			if err == nil {
				if test.err != nil {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			} else {
				if err.Error() != test.err.Error() {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			}
			if !reflect.DeepEqual(res, test.res) {
				t.Fatalf("expected response to be %q but got %q", test.res, res)
			}
		})
	}
}

// TestDecodeYAML tests the DecodeYAML function
func TestDecodeYAML(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input io.Reader
		name  string
		res   interface{}
		err   error
	}{
		{
			input: strings.NewReader("value: -"),
			name:  "decode-yaml-error",
			err:   errors.New("yaml: block sequence entries are not allowed in this context"),
		},
		{
			input: strings.NewReader("a: b"),
			name:  "decode-yaml-success",
			res:   map[string]interface{}{"a": "b"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeYAML(test.input)
			if err == nil {
				if test.err != nil {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			} else {
				if err.Error() != test.err.Error() {
					t.Fatalf("expected error to be %q but got %q", test.err, err)
				}
			}
			//			if !reflect.DeepEqual(res, test.res) {
			//				t.Fatalf("expected response to be %q but got %q", test.res, res)
			//			}
		})
	}
}

const (
	errorDecodeYAMLTest = `apiVersion: v1
kind: Secret
metadata:
  name: error-test4
type: Opaque
data:
  a: 世界
`

	successDecodeJSONTestEmpty = `{
	"apiVersion": "v1",
	"data": {},
	"kind": "Secret",
	"metadata": {
		"name": "json-empty"
	},
	"type": "Opaque"
}`

	successDecodeJSONTestOne = `{
	"apiVersion": "v1",
	"data": {
		"a": "Yg=="
	},
	"kind": "Secret",
	"metadata": {
		"name": "json-one"
	},
	"type": "Opaque"
}`

	successDecodeYAMLTestEmpty = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-empty
type: Opaque
data:
`

	successDecodeYAMLTestEmpty2 = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-empty2
type: Opaque
data: {}
`

	successDecodeYAMLTestOne = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-one
type: Opaque
data:
  a: Yg==
`
)
