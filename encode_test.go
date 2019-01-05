package k8shhh

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// TestEncode tests the Encode function
func TestEncode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   io.Reader
		encoder Encoder
		name    string
		res     string
		err     error
	}{
		{
			input:   strings.NewReader("-1"),
			encoder: EncodeYAML,
			name:    "error-test",
			err:     errors.New("Can't separate key from value"),
		},
		{
			input:   strings.NewReader(""),
			encoder: EncodeJSON,
			name:    "json-empty",
			res:     successEncodeJSONTestEmpty,
		},
		{
			input:   strings.NewReader(""),
			encoder: EncodeYAML,
			name:    "yaml-empty",
			res:     successEncodeYAMLTestEmpty,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeJSON,
			name:    "json-one",
			res:     successEncodeJSONTestOne,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeYAML,
			name:    "yaml-one",
			res:     successEncodeYAMLTestOne,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := Encode(test.input, test.encoder, test.name)
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

// TestEncodeJSON tests the EncodeJSON function
func TestEncodeJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"json-empty", make(map[string]string)},
			res:    successEncodeJSONTestEmpty,
		},
		{
			secret: Secret{"json-one", map[string]string{"a": "b"}},
			res:    successEncodeJSONTestOne,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.secret.Name, func(t *testing.T) {
			res, err := EncodeJSON(test.secret)
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

// TestEncodeYAML tests the EncodeYAML function
func TestEncodeYAML(t *testing.T) {
	t.Parallel()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"yaml-empty", make(map[string]string)},
			res:    successEncodeYAMLTestEmpty,
		},
		{
			secret: Secret{"yaml-one", map[string]string{"a": "b"}},
			res:    successEncodeYAMLTestOne,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.secret.Name, func(t *testing.T) {
			res, err := EncodeYAML(test.secret)
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

const (
	successEncodeJSONTestEmpty = `{
	"apiVersion": "v1",
	"data": {},
	"kind": "Secret",
	"metadata": {
		"name": "json-empty"
	},
	"type": "Opaque"
}`

	successEncodeJSONTestOne = `{
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

	successEncodeYAMLTestEmpty = `apiVersion: v1
data: {}
kind: Secret
metadata:
  name: yaml-empty
type: Opaque
`

	successEncodeYAMLTestOne = `apiVersion: v1
data:
  a: Yg==
kind: Secret
metadata:
  name: yaml-one
type: Opaque
`
)
