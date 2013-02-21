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
	fmt.Fprintf(os.Stderr, "Usage: jsonget file.json attribute[.subattribute] [...]\n")
	os.Exit(2)
}

func main() {
	flag.Parse()

	filePath := flag.Arg(0)
	properties := flag.Args()[1:]

	if len(filePath) > 0 && len(properties) > 0 {
		data, err := jsonFromFile(filePath)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		values, err := getValues(&data, properties)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
		for _, value := range values {
			fmt.Println(value)
		}
	} else {
		usage()
	}
}
