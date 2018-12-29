package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/jwangsadinata/k8shhh"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// VERSION is the k8shhh app version
	VERSION = "1.0.0"
)

var (
	app = kingpin.New("k8shhh", "k8shhh: Quickly encode your configuration into K8s secrets.")

	enc           = app.Command("encode", "encode your configuration as k8s secrets")
	encSecretName = enc.Flag("name", "the name of the generated secret").Short('n').String()
	encInput      = enc.Flag("input", "the name of the input file to encode (if input is not provided via STDIN)").Short('i').String()
	encOutput     = enc.Flag("output", "the name of the file to write the output to (outputs to STDOUT by default). file extension will be automatically generated based on the format.").Short('o').String()
	encFormat     = enc.Flag("format", "format of the generated secret (json or yaml, defaults to yaml)").Default("yaml").Short('f').String()

	dec       = app.Command("decode", "decode your k8s secrets into a readable format")
	decInput  = dec.Flag("input", "the name of the input file to decode (if input is not provided via STDIN)").Short('i').String()
	decOutput = dec.Flag("output", "the name of the file to write the output to (outputs to STDOUT by default)").Short('o').String()

	version = app.Command("version", "print the current version of k8shhh.")
)

// run initializes the command line parser
func run() int {
	kingpin.Version(VERSION)

	ctx, err := app.ParseContext(os.Args[1:])
	if err != nil {
		return 1
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case enc.FullCommand():
		if isInteractive() && *encInput == "" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "expecting input on stdin")
			return 1
		}

		if *encFormat != "json" && *encFormat != "yaml" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "format must be either yaml or json")
			return 1
		}

		var input io.Reader
		input = os.Stdin
		if *encInput != "" {
			f, err := os.Open(*encInput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reading input file: %s", err)
				return 1
			}
			defer f.Close()
			input = f
		}

		var encoder Encoder
		switch *encFormat {
		case "json":
			encoder = EncodeJson
		case "yaml":
			encoder = EncodeYaml
		}

		var secretName string
		if *encSecretName == "" {
			if *encOutput == "" {
				secretName = "mysecret"
			} else {
				secretName = trimExtension(*encOutput)
			}
		} else {
			secretName = *encSecretName
		}

		if output, err := Encode(input, encoder, secretName); err != nil {
			fmt.Fprintf(os.Stderr, "error in encoding: %v\n", err)
			return 1
		} else {
			if *encOutput != "" {
				filename := *encOutput
				if !checkExtension(*encOutput) {
					filename = fmt.Sprintf("%s.%s", *encOutput, *encFormat)
				}
				err := ioutil.WriteFile(filename, output, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "writing to output file: %s", err)
					return 1
				}

				// print the file name out - useful for pipe with `kubectl create -f`
				fmt.Print(filename)
			} else {
				fmt.Print(string(output))
			}
		}
	case dec.FullCommand():
		if isInteractive() && *decInput == "" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "expecting input on stdin")
			return 1
		}

		var input io.Reader
		input = os.Stdin
		if *decInput != "" {
			f, err := os.Open(*decInput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reading input file: %s", err)
				return 1
			}
			defer f.Close()
			input = f
		}

		var decoder Decoder
		if *decInput != "" && strings.Contains(*decInput, ".json") {
			decoder = DecodeJson
		} else {
			decoder = DecodeYaml
		}

		if output, err := Decode(input, decoder); err != nil {
			fmt.Fprintf(os.Stderr, "error in decoding: %v\n", err)
			return 1
		} else {
			if *decOutput != "" {
				err := ioutil.WriteFile(*decOutput, output, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "writing to output file: %s", err)
					return 1
				}

				fmt.Printf(`file "%s" created`, *decOutput)
			} else {
				fmt.Print(string(output))
			}
		}
	case version.FullCommand():
		fmt.Printf("k8shhh %s\n", VERSION)
	}

	return 0
}

// checkExtension checks whether the output string contains a .json or .yaml
// extension
func checkExtension(output string) bool {
	return strings.HasSuffix(output, ".yaml") || strings.HasSuffix(output, ".json")
}

// trimExtension trims the output string if it contains a .json or .yaml
// extension
func trimExtension(output string) string {
	if checkExtension(output) {
		return output[:len(output)-5]
	}
	return output
}

// isInteractive returns true if os.Stdin appears to be interactive
func isInteractive() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&(os.ModeCharDevice|os.ModeCharDevice) != 0
}

// main executes the run function
func main() {
	os.Exit(run())
}
