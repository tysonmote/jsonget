package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	filePath   = flag.String("file", "", "Read from a file instead of stdin.")
	printNulls = flag.Bool("nulls", false, "If true, null values will be printed as 'null'.")
	silent     = flag.Bool("silent", false, "If true, errors will not be printed to stderr.")
)

func usage() {
	fmt.Printf(`Usage: jsonget -file [JSON_FILE] attribute ... [attribute]

Examples:
  cat data.json | jsonget person.name person.age
  jsonget -file data.json person.address.city`)
	os.Exit(1)
}

// Die if given an error.
func dieIfError(err error) {
	if err != nil {
		if !*silent {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(1)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Load command-line attribute chains

	attributes := flag.Args()
	if len(attributes) == 0 {
		usage()
	}

	// Read in the JSON

	var data *JsonData
	var err error

	if len(*filePath) > 0 {
		data, err = LoadFile(*filePath)
	} else {
		data, err = LoadStdin()
	}
	dieIfError(err)

	// Get and print the values from the JSON

	for _, attribute := range attributes {
		values, err := data.GetValues(attribute)
		dieIfError(err)

		for _, value := range values {
			fmt.Println(value)
		}
	}
}
