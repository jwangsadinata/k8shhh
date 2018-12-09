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
			encoder: EncodeYaml,
			name:    "error-test",
			err:     errors.New("Can't separate key from value"),
		},
		{
			input:   strings.NewReader(""),
			encoder: EncodeJson,
			name:    "json-empty",
			res:     successEncodeJsonTestEmpty,
		},
		{
			input:   strings.NewReader(""),
			encoder: EncodeYaml,
			name:    "yaml-empty",
			res:     successEncodeYamlTestEmpty,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeJson,
			name:    "json-one",
			res:     successEncodeJsonTestOne,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeYaml,
			name:    "yaml-one",
			res:     successEncodeYamlTestOne,
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

// TestEncodeJson tests the EncodeJson function
func TestEncodeJson(t *testing.T) {
	t.Parallel()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"json-empty", make(map[string]string)},
			res:    successEncodeJsonTestEmpty,
		},
		{
			secret: Secret{"json-one", map[string]string{"a": "b"}},
			res:    successEncodeJsonTestOne,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.secret.Name, func(t *testing.T) {
			res, err := EncodeJson(test.secret)
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

// TestEncodeYaml tests the EncodeYaml function
func TestEncodeYaml(t *testing.T) {
	t.Parallel()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"yaml-empty", make(map[string]string)},
			res:    successEncodeYamlTestEmpty,
		},
		{
			secret: Secret{"yaml-one", map[string]string{"a": "b"}},
			res:    successEncodeYamlTestOne,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.secret.Name, func(t *testing.T) {
			res, err := EncodeYaml(test.secret)
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
	successEncodeJsonTestEmpty = `{
	"apiVersion": "v1",
	"data": {},
	"kind": "Secret",
	"metadata": {
		"name": "json-empty"
	},
	"type": "Opaque"
}`

	successEncodeJsonTestOne = `{
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

	successEncodeYamlTestEmpty = `apiVersion: v1
data: {}
kind: Secret
metadata:
  name: yaml-empty
type: Opaque
`

	successEncodeYamlTestOne = `apiVersion: v1
data:
  a: Yg==
kind: Secret
metadata:
  name: yaml-one
type: Opaque
`
)
