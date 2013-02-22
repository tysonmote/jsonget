package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	quotedString = regexp.MustCompile("\\A\"(.+)\"\\z")
)

type JsonObject map[string]interface{}

func (data JsonObject) GetValue(attribute string) (value string, err error) {
	attributeParts := strings.Split(attribute, ".")
	attributePartsCount := len(attributeParts)

	var cursor JsonObject
	cursor = data

	for i, attributePart := range(attributeParts) {
		nextCursor := cursor[attributePart]

		if i == attributePartsCount - 1 || nextCursor == nil {
			// Last attribute part
			return valueToString(nextCursor)
		} else if isMap(nextCursor) {
			cursor = JsonObject(nextCursor.(map[string]interface{}))
		} else {
			parentAttribute := strings.Join(attributeParts[0:i+1], ".")
			err := fmt.Errorf("Can't read %s attribute on %s because it is not a JSON object.", attributeParts[i+1], parentAttribute )
			return "", err
		}
	}

	return valueToString(cursor)
}

func (data JsonObject) GetValues(attributeChains []string) (values []string, err error) {
	values = make([]string, len(attributeChains))

	for i, attributeChain := range attributeChains {
		value, err := data.GetValue(attributeChain)
		if err != nil {
			return values, err
		}
		values[i] = value
	}

	return values, nil
}

// unmarshal takes a byte array and parses it into a JsonObject structure.
func unmarshal(text []byte) (jsonData JsonObject, err error) {
	var data JsonObject
	err = json.Unmarshal(text, &data)
	if err != nil {
		return JsonObject{}, err
	}

	return data, nil
}

// JsonObjectFromFile reads JSON data from the given file path and parses it into
// a JsonObject object.
func JsonObjectFromFile(file string) (jsonData JsonObject, err error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return make(JsonObject), err
	}
	return unmarshal(text)
}

// JsonObjectFromStdin reads stdin for JSON data and parses it into a JsonObject object.
func JsonObjectFromStdin() (jsonData JsonObject, err error) {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return make(JsonObject), err
	}
	return unmarshal(text)
}

// valueToString returns a string representation of the given value. If the
// given value is nil, a blank string is returned. Or, if the printNulls flag
// is true, a "null" string is returned.
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

