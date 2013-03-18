package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
	"regexp"
	"strings"
	"strconv"
)

var (
	quotedString = regexp.MustCompile("\\A\"(.+)\"\\z")
)

type JsonObject map[string]interface{}

// Find and return the given attribute's value. The attribute can use
// dot-notation ("foo.bar.baz") to access inner attributes. If the value is a
// string, it is returned without quote marks. Otherwise, it is returned as a
// JSON string.
func (data JsonObject) GetValue(attribute string) (value string, err error) {
	attributeParts := strings.Split(attribute, ".")
	attributePartsCount := len(attributeParts)

	var cursor JsonObject
	cursor = data

	for i, attributePart := range attributeParts {
		var nextCursor interface{}

		// if attributePart has '[<int>]' in it, then index into array
		aryExp, aryIndices, attrPartTrimmed := isArrayExpression(attributePart)

		if aryExp {
			nextCursor, _ = cursor[attrPartTrimmed].([]interface{})

			for i:=0; i < len(aryIndices); i++ {
				arrayCursor, _ := nextCursor.([]interface{})

				if aryIndices[i] < uint64(len(arrayCursor)) {
					nextCursor = arrayCursor[aryIndices[i]]
				} else {
					return "", errors.New("Index out of bounds of JSON array")
				}
			}
		} else {
			nextCursor = cursor[attributePart]
		}

		if i == attributePartsCount-1 || nextCursor == nil {
			return valueToString(nextCursor)
		} else {
			nextCursorMap, ok := nextCursor.(map[string]interface{})
			if ok {
				cursor = JsonObject(nextCursorMap)
			} else {
				parentAttribute := strings.Join(attributeParts[0:i+1], ".")
				err := fmt.Errorf("Can't read %s attribute on %s because it is not a JSON object.", attributeParts[i+1], parentAttribute)
				return "", err
			}
		}
	}

	return valueToString(cursor)
}

// Get several values at once. See GetValue for attribute string formatting
// rules.
func (data JsonObject) GetValues(attributes []string) (values []string, err error) {
	values = make([]string, len(attributes))

	for i, attribute := range attributes {
		value, err := data.GetValue(attribute)
		if err != nil {
			return values, err
		}
		values[i] = value
	}

	return values, nil
}

// Parse the given JSON into a JsonObject object.
func unmarshal(text []byte) (jsonData JsonObject, err error) {
	var data JsonObject
	if err = json.Unmarshal(text, &data); err != nil {
		return JsonObject{}, err
	}

	return data, nil
}

// Read the given JSON data and parse it into a JsonObject object;
func JsonObjectFromFile(file string) (jsonData JsonObject, err error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return make(JsonObject), err
	}
	return unmarshal(text)
}

// Read stdin for JSON and parse it into a JsonObject object.
func JsonObjectFromStdin() (jsonData JsonObject, err error) {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return make(JsonObject), err
	}
	return unmarshal(text)
}

// Get a string representation of the given value. If the given value is nil, a
// blank string is returned. Or, if the printNulls flag is true, a "null"
// string is returned.
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

// Check if attribute is an array expression (e.g.  whoa[3], stuff[1][5], ... etc)
//   Return: - whether it is array expression
//           - list of index values
//           - name of json key
func isArrayExpression(attribute string) (bool, []uint64, string) {
	reArray, _ := regexp.Compile(`\[\d+\]`)

	//                                    magic ↴
	matches := reArray.FindAllString(attribute, 42)
	n := len(matches)
	indices := []uint64{}
	if n > 0 {
		for i:=0; i < n; i++ {
			sIndex := strings.Replace(matches[i][1:], `]`, "", 1)
			index, _ := strconv.ParseUint(sIndex, 10, 64)
			indices = append(indices, index)
		}
		bidx := strings.Index(attribute, `]`)
		if bidx != -1 {
			attribute = attribute[:bidx-2]
		}
	}

		return n > 0, indices, attribute
	}

