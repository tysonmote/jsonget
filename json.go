package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var quotedString = regexp.MustCompile(`\A"(.+)"\z`)

type JsonData struct {
	json interface{}
}

// Read the given JSON data and parse it into a JsonObject object;
func LoadFile(file string) (data *JsonData, err error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	object, err := unmarshal(text)
	if err != nil {
		return nil, err
	}
	return &JsonData{object}, nil
}

// Read stdin for JSON and parse it into a JsonObject object.
func LoadStdin() (data *JsonData, err error) {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	object, err := unmarshal(text)
	if err != nil {
		return nil, err
	}
	return &JsonData{object}, nil
}

// Find and return the given attribute's value. The attribute can use
// dot-notation ("foo.bar.baz") to access inner attributes. If the value is a
// string, it is returned without quote marks. Otherwise, it is returned as a
// JSON string.
func (j *JsonData) GetValue(attribute string) (value string, err error) {
	attributeParts := splitAttributeParts(attribute)
	attributePartsCount := len(attributeParts)

	cursor := j.json

	for i, attributePart := range attributeParts {
		var nextCursor interface{}

		if cursor == nil {
			errorString := fmt.Sprintf("Cannot access \"%s\" on nil.", attributePart)
			return "", errors.New(errorString)
		}

		cursorKind := reflect.ValueOf(cursor).Kind()

		switch cursorKind {
		case reflect.Map:
			nextCursor = cursor.(map[string]interface{})[attributePart]
		case reflect.Slice:
			index, err := strconv.ParseInt(attributePart, 10, 0)
			if err != nil {
				parentAttribute := strings.Join(attributeParts[0:i], ".")
				errorString := fmt.Sprintf("%s is an array, but %s is not a valid index.", parentAttribute, attributePart)
				return "", errors.New(errorString)
			}
			cursorSlice := cursor.([]interface{})
			if int(index) >= len(cursorSlice) {
				parentAttribute := strings.Join(attributeParts[0:i], ".")
				errorString := fmt.Sprintf("%d is outside the bounds of the %d elements in %s.", index, len(cursorSlice), parentAttribute)
				return "", errors.New(errorString)
			}
			nextCursor = cursorSlice[index]
		default:
			parentAttribute := strings.Join(attributeParts[0:i], ".")
			errorString := fmt.Sprintf("Can't get %s on %s because it is a %v.", attributePart, parentAttribute, cursorKind)
			return "", errors.New(errorString)
		}

		if i == attributePartsCount-1 || nextCursor == nil {
			return valueToString(nextCursor)
		} else {
			cursor = nextCursor
		}
	}

	return valueToString(cursor)
}

// Get several values at once. See GetValue for attribute string formatting
// rules.
func (j *JsonData) GetValues(attributes []string) (values []string, err error) {
	values = make([]string, len(attributes))

	for i, attribute := range attributes {
		value, err := j.GetValue(attribute)
		if err != nil {
			return values, err
		}
		values[i] = value
	}

	return values, nil
}

// Parse the given JSON.
func unmarshal(text []byte) (object interface{}, err error) {
	var jsonData interface{}
	err = json.Unmarshal(text, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
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

// Returns all the attribute parts, including turning array access into "plain"
// attribute access. This is something of a hack.
//
// Examples:
//   * "foo.bar" --> []string{"foo", "bar"}
//   * "foo.bar[2].neat --> [string]{"foo", "bar", "2", "neat"}
func splitAttributeParts(attribute string) []string {
	brackets := regexp.MustCompile(`[\[\]]+`)
	dots := regexp.MustCompile(`\.+`)

	attributeBytes := brackets.ReplaceAll([]byte(attribute), []byte{'.'})
	attributeBytes = dots.ReplaceAll(attributeBytes, []byte{'.'})

	if bytes.LastIndex(attributeBytes, []byte{'.'}) == len(attributeBytes)-1 {
		attributeBytes = attributeBytes[:len(attributeBytes)-1]
	}
	return strings.Split(string(attributeBytes), ".")
}
