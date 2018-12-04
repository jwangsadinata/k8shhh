package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	. "github.com/jwangsadinata/k8shhh"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	VERSION = "0.0.1"
)

var (
	app = kingpin.New("k8shhh", "k8shhh: Quickly encode your configuration into K8s secrets.")

	enc           = app.Command("encode", "encode your configuration as k8s secrets")
	encSecretName = enc.Flag("name", "the name of the generated secret").Default("mysecret").Short('n').String()
	encInput      = enc.Flag("input", "the name of the input file to encode (if input is not provided via STDIN)").Short('i').String()
	encOutput     = enc.Flag("output", "the name of the file to write the output to (outputs to STDOUT by default)").Short('o').String()
	encFormat     = enc.Flag("format", "format of the generated secret (json or yaml, defaults to yaml)").Default("yaml").Short('f').String()

	version = app.Command("version", "print the current version of k8shhh.")
)

// run initializes the command line parser
func run() int {
	// set the version
	kingpin.Version(VERSION)

	// parse the context for the application
	ctx, err := app.ParseContext(os.Args[1:])
	if err != nil {
		return 1
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case enc.FullCommand():
		// validate encoding input
		if isInteractive() && *encInput == "" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "expecting input on stdin")
			return 1
		}

		// validate encoding format
		if *encFormat != "json" && *encFormat != "yaml" {
			kingpin.CommandLine.UsageForContext(ctx)
			fmt.Fprintln(os.Stderr, "format must be either yaml or json")
			return 1
		}

		// setup the encoding input
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

		// parse the encoder
		var encoder Encoder
		switch *encFormat {
		case "json":
			encoder = EncodeJson
		case "yaml":
			encoder = EncodeYaml
		}

		// encode the secret file
		if output, err := Encode(input, encoder, *encSecretName); err != nil {
			fmt.Fprintf(os.Stderr, "error in encoding: %v\n", err)
			return 1
		} else {
			if *encOutput != "" {
				// write to file with the appropriate permissions
				err := ioutil.WriteFile(*encOutput, output, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "writing to output file: %s", err)
					return 1
				}

				// print the file name out - useful for kubectl create secrets
				fmt.Print(*encOutput)
			} else {
				// print the output to stdout
				fmt.Print(string(output))
			}
		}
	case version.FullCommand():
		// print out the version
		fmt.Printf("k8shhh %s\n", VERSION)
	}

	return 0
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
