package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: json file.json attribute[.subattribute] [...]\n")
	os.Exit(2)
}

func jsonFromFile(file string) (jsonData map[string]interface{}, err error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return make(map[string]interface{}), err
	}

	var data interface{}
	err = json.Unmarshal(text, &data)
	if err != nil {
		return make(map[string]interface{}), err
	}

	return data.(map[string]interface{}), nil
}

func get(data *map[string]interface{}, attributeChain []string) (value string, err error) {
	attribute := attributeChain[0]

	if len(attributeChain) == 1 {
		text, err := json.Marshal((*data)[attribute])
		return string(text), err
	}

	subdata := (*data)[attribute].(map[string]interface{})
	return get(&subdata, attributeChain[1:])
}

func getValues(data *map[string]interface{}, attributeChains []string) (values []string, err error) {
	values = make([]string, len(attributeChains))
	for i, attributeChain := range(attributeChains) {
		value, err := get(data, strings.Split(attributeChain, "."))
		if err != nil {
			return values, nil
		}
		values[i] = value
	}
	return values, nil
}

func main() {
	flag.Parse()

	path := flag.Arg(0)
	properties := flag.Args()[1:]

	if len(path) > 0 && len(properties) > 0 {
		data, err := jsonFromFile(path)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		values, err := getValues(&data, properties)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
		for _, value := range(values) {
			fmt.Println(value)
		}
	} else {
		usage()
	}
}
