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
			decoder: DecodeYaml,
			name:    "error-test",
			err:     errors.New("unexpected type: int"),
		},
		{
			input:   strings.NewReader("data: -1"),
			decoder: DecodeYaml,
			name:    "error-test2",
			err:     errors.New("unexpected type: int"),
		},
		{
			input:   strings.NewReader("value: -"),
			decoder: DecodeYaml,
			name:    "error-test3",
			err:     errors.New("yaml: block sequence entries are not allowed in this context"),
		},
		{
			input:   strings.NewReader(errorDecodeYamlTest),
			decoder: DecodeYaml,
			name:    "error-test4",
			err:     errors.New("illegal base64 data at input byte 0"),
		},
		{
			input:   strings.NewReader(successDecodeJsonTestEmpty),
			decoder: DecodeJson,
			name:    "json-empty",
			res:     "",
		},
		{
			input:   strings.NewReader(successDecodeYamlTestEmpty),
			decoder: DecodeYaml,
			name:    "yaml-empty",
			res:     "",
		},
		{
			input:   strings.NewReader(successDecodeJsonTestOne),
			decoder: DecodeJson,
			name:    "json-one",
			res:     "a=b",
		},
		{
			input:   strings.NewReader(successDecodeYamlTestOne),
			decoder: DecodeYaml,
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

// TestDecodeJson tests the DecodeJson function
func TestDecodeJson(t *testing.T) {
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
			res, err := DecodeJson(test.input)
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

// TestDecodeYaml tests the DecodeYaml function
func TestDecodeYaml(t *testing.T) {
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
			_, err := DecodeYaml(test.input)
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
	errorDecodeYamlTest = `apiVersion: v1
kind: Secret
metadata:
  name: error-test4
type: Opaque
data:
  a: 世界
`

	successDecodeJsonTestEmpty = `{
	"apiVersion": "v1",
	"data": {},
	"kind": "Secret",
	"metadata": {
		"name": "json-empty"
	},
	"type": "Opaque"
}`

	successDecodeJsonTestOne = `{
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

	successDecodeYamlTestEmpty = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-empty
type: Opaque
data:
`

	successDecodeYamlTestOne = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-one
type: Opaque
data:
  a: Yg==
`
)
