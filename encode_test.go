package k8shhh

import (
	"io"
	"strings"
	"testing"
)

var dataMap = make(map[string]string)

// initializeDataMap initializes the data map
func initializeDataMap() {
	dataMap["a"] = "b"
}

// TestEncode tests the Encode function
func TestEncode(t *testing.T) {
	tests := []struct {
		input   io.Reader
		encoder Encoder
		name    string
		res     string
		err     error
	}{
		{
			input:   strings.NewReader(""),
			encoder: EncodeJson,
			name:    "json-empty",
			res:     successJsonTestEmpty,
			err:     nil,
		},
		{
			input:   strings.NewReader(""),
			encoder: EncodeYaml,
			name:    "yaml-empty",
			res:     successYamlTestEmpty,
			err:     nil,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeJson,
			name:    "json-one",
			res:     successJsonTestOne,
			err:     nil,
		},
		{
			input:   strings.NewReader("a=b"),
			encoder: EncodeYaml,
			name:    "yaml-one",
			res:     successYamlTestOne,
			err:     nil,
		},
	}

	for _, test := range tests {
		res, err := Encode(test.input, test.encoder, test.name)
		if err != test.err {
			t.Fatalf("expected error to be %q but got %q", test.err, err)
		}
		if string(res) != test.res {
			t.Fatalf("expected response to be %q but got %q", test.res, res)
		}
	}
}

// TestEncodeJson tests the EncodeJson function
func TestEncodeJson(t *testing.T) {
	initializeDataMap()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"json-empty", make(map[string]string)},
			res:    successJsonTestEmpty,
			err:    nil,
		},
		{
			secret: Secret{"json-one", dataMap},
			res:    successJsonTestOne,
			err:    nil,
		},
	}

	for _, test := range tests {
		res, err := EncodeJson(test.secret)
		if err != test.err {
			t.Fatalf("expected error to be %q but got %q", test.err, err)
		}
		if string(res) != test.res {
			t.Fatalf("expected response to be %q but got %q", test.res, res)
		}
	}
}

// TestEncodeYaml tests the EncodeYaml function
func TestEncodeYaml(t *testing.T) {
	initializeDataMap()
	tests := []struct {
		secret Secret
		res    string
		err    error
	}{
		{
			secret: Secret{"yaml-empty", make(map[string]string)},
			res:    successYamlTestEmpty,
			err:    nil,
		},
		{
			secret: Secret{"yaml-one", dataMap},
			res:    successYamlTestOne,
			err:    nil,
		},
	}

	for _, test := range tests {
		res, err := EncodeYaml(test.secret)
		if err != test.err {
			t.Fatalf("expected error to be %q but got %q", test.err, err)
		}
		if string(res) != test.res {
			t.Fatalf("expected response to be %q but got %q", test.res, res)
		}
	}
}

const (
	successJsonTestEmpty = `{
	"apiVersion": "v1",
	"data": {},
	"kind": "Secret",
	"metadata": {
		"name": "json-empty"
	},
	"type": "Opaque"
}`

	successJsonTestOne = `{
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

	successYamlTestEmpty = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-empty
type: Opaque
data:
`

	successYamlTestOne = `apiVersion: v1
kind: Secret
metadata:
  name: yaml-one
type: Opaque
data:
  a: Yg==
`
)
