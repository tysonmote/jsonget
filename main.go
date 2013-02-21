package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	printNulls = flag.Bool("nulls", false, "If true, null values will be printed as 'null'.")
)

func usage() {
	die("Usage: jsonget file.json attribute[.subattribute] [...]")
}

func die(text string) {
	fmt.Fprintf(os.Stderr, "%s\n", text)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	filePath := flag.Arg(0)
	properties := flag.Args()[1:]

	if len(filePath) > 0 && len(properties) > 0 {
		data, err := jsonFromFile(filePath)
		if err != nil {
			die(err.Error())
		}

		values, err := getValues(&data, properties)
		if err != nil {
			die(err.Error())
		}
		for _, value := range values {
			fmt.Println(value)
		}
	} else {
		usage()
	}
}
