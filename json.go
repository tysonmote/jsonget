package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	quotedString = regexp.MustCompile("\\A\"(.+)\"\\z")
)

type JsonData map[string]interface{}

func (data JsonData) GetValue(attributeChain []string) (value string, err error) {
	attribute := attributeChain[0]

	if len(attributeChain) == 1 {
		return valueToString(data[attribute])
	}

	rawSubdata := data[attribute]
	if rawSubdata == nil {
		return "", nil
	}

	subdata := JsonData(rawSubdata.(map[string]interface{}))
	return subdata.GetValue(attributeChain[1:])
}

func (data JsonData) GetValues(attributeChains []string) (values []string, err error) {
	values = make([]string, len(attributeChains))

	for i, attributeChain := range attributeChains {
		value, err := data.GetValue(strings.Split(attributeChain, "."))
		if err != nil {
			return values, err
		}
		values[i] = value
	}

	return values, nil
}

func unmarshal(text []byte) (jsonData JsonData, err error) {
	var data JsonData
	err = json.Unmarshal(text, &data)
	if err != nil {
		return JsonData{}, err
	}

	return data, nil
}

func JsonDataFromFile(file string) (jsonData JsonData, err error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return make(JsonData), err
	}
	return unmarshal(text)
}

func JsonDataFromStdin() (jsonData JsonData, err error) {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return make(JsonData), err
	}
	return unmarshal(text)
}

func valueToString(value interface{}) (text string, err error) {
	if value == nil && *printNulls == false {
		return "", nil
	}

	textBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	text = string(textBytes)
	text = quotedString.ReplaceAllString(text, "$1")
	return text, nil
}

