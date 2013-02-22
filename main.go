package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	filePath = flag.String("file", "", "Read from a file instead of stdin.")
	printNulls = flag.Bool("nulls", false, "If true, null values will be printed as 'null'.")
)

func usage() {
	die("Usage: jsonget -f file.json attribute[.subattribute] [...]")
}

func die(text string) {
	fmt.Fprintf(os.Stderr, "%s\n", text)
	os.Exit(1)
}

func dieIfError(err error) {
	if err != nil {
		die(err.Error())
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	properties := flag.Args()

	if len(properties) == 0 {
		usage()
	}

	var data JsonData
	var err error

	if len(*filePath) > 0 {
		data, err = JsonDataFromFile(*filePath)
	} else {
		data, err = JsonDataFromStdin()
	}
	dieIfError(err)

	values, err := data.GetValues(properties)
	dieIfError(err)

	for _, value := range values {
		fmt.Println(value)
	}
}
