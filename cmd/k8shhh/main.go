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

		input, err := selectInput(*encInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading input file: %s", err)
			return 1
		}
		defer input.Close()

		encoder := selectEncoder(*encFormat)
		secretName := initializeSecretName(*encSecretName, *encOutput)
		output, err := Encode(input, encoder, secretName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in encoding: %v\n", err)
			return 1
		}

		msg, err := processEncodeOutput(output, *encOutput, *encFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "writing to output file: %s", err)
			return 1
		}
		fmt.Print(msg)
	case dec.FullCommand():
		if isInteractive() && *decInput == "" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "expecting input on stdin")
			return 1
		}

		input, err := selectInput(*decInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading input file: %s", err)
			return 1
		}

		decoder := selectDecoder(*decInput)
		output, err := Decode(input, decoder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in decoding: %v\n", err)
			return 1
		}

		msg, err := processDecodeOutput(output, *decOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "writing to output file: %s", err)
			return 1
		}
		fmt.Print(msg)
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

// checkFormat checks whether the format passed is correct
func checkFormat(format string) int {
	return 0
}

// initializeSecretName initializes the secret name
func initializeSecretName(sn, output string) string {
	if sn == "" {
		if output == "" {
			return "mysecret"
		}
		return trimExtension(*encOutput)
	}
	return sn
}

// processDecodeOutput process the output of the decoder
func processDecodeOutput(output []byte, file string) (string, error) {
	if file != "" {
		err := ioutil.WriteFile(file, output, 0644)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`file "%s" created`, *decOutput), nil
	}
	return string(output), nil
}

// processEncodeOutput process the output of the encoder
func processEncodeOutput(output []byte, file, format string) (string, error) {
	if file != "" {
		filename := file
		if !checkExtension(file) {
			filename = fmt.Sprintf("%s.%s", file, format)
		}
		err := ioutil.WriteFile(filename, output, 0644)
		if err != nil {
			return "", err
		}
		// print the file name out - useful for pipe with `kubectl create -f`
		return filename, nil
	}
	return string(output), nil
}

// selectDecoder returns an encoder based on the input provided.
func selectDecoder(input string) Decoder {
	if input != "" && strings.Contains(input, ".json") {
		return DecodeJSON
	}
	return DecodeYAML
}

// selectEncoder returns an encoder based on the format provided.
func selectEncoder(format string) Encoder {
	if format == "json" {
		return EncodeJSON
	}
	return EncodeYAML
}

// selectInput returns the io.Reader based on the provided input.
func selectInput(s string) (io.ReadCloser, error) {
	var input io.ReadCloser
	input = os.Stdin
	if s != "" {
		f, err := os.Open(s)
		if err != nil {
			return nil, err
		}
		input = f
	}
	return input, nil
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
	return fileInfo.Mode()&os.ModeCharDevice != 0
}

// main executes the run function
func main() {
	os.Exit(run())
}
