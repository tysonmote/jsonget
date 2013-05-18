package main

import (
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

var quotedStringRegex = regexp.MustCompile(`\A"(.+)"\z`)
var partsRegex = regexp.MustCompile(`[\[\]\.]+`)

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
	return &JsonData{json: object}, nil
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
	return &JsonData{json: object}, nil
}

func (j *JsonData) GetValues(attribute string) (values []string, err error) {
	vals, err := recursiveGetValues(0, splitAttributeParts(attribute), []interface{}{j.json})
	if err != nil {
		return []string{}, err
	}

	return valuesToStrings(vals)
}

func recursiveGetValues(depth int, attributeParts []string, objects []interface{}) (values []interface{}, err error) {
	mappedData := []interface{}{}

	if depth == len(attributeParts) {
		return objects, nil
	}

	part := attributeParts[depth]

	for _, object := range objects {
		if object == nil {
			errorString := fmt.Sprintf(`Cannot access "%s" on nil`, part)
			return []interface{}{}, errors.New(errorString)
		}

		objectKind := reflect.TypeOf(object).Kind()

		switch objectKind {
		case reflect.Map:
			mappedData = append(mappedData, object.(map[string]interface{})[part])
		case reflect.Slice:
			if part == "*" {
				mappedData = append(mappedData, object.([]interface{})...)
			} else {
				index, err := strconv.ParseInt(part, 10, 0)
				if err != nil {
					parentAttribute := strings.Join(attributeParts[0:depth], ".")
					errorString := fmt.Sprintf(`"%s" is not a valid index for the array at "%s"`, part, parentAttribute)
					return []interface{}{}, errors.New(errorString)
				}
				slice := object.([]interface{})
				if int(index) >= len(slice) {
					parentAttribute := strings.Join(attributeParts[0:depth], ".")
					errorString := fmt.Sprintf(`%d is outside the bounds of the %d elements in "%s"`, index, len(slice), parentAttribute)
					return []interface{}{}, errors.New(errorString)
				}
				mappedData = append(mappedData, slice[index])
			}
		default:
			parentAttribute := strings.Join(attributeParts[0:depth], ".")
			errorString := fmt.Sprintf(`Can't get "%s" on "%s" because it is a %v`, part, parentAttribute, objectKind)
			return []interface{}{}, errors.New(errorString)
		}
	}

	return recursiveGetValues(depth+1, attributeParts, mappedData)
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
	text = quotedStringRegex.ReplaceAllString(text, "$1")

	return text, nil
}

func valuesToStrings(values []interface{}) (strings []string, err error) {
	for _, value := range values {
		str, err := valueToString(value)
		if err != nil {
			return []string{}, err
		}
		strings = append(strings, str)
	}
	return strings, nil
}

// Returns all the attribute parts, including turning array access into "plain"
// attribute access. This is something of a hack.
//
// Examples:
//   * "foo.bar" --> []string{"foo", "bar"}
//   * "foo.bar[2].neat --> []string{"foo", "bar", "2", "neat"}
//   * "cities.*.name --> []string{"cities", "*", "name"}
func splitAttributeParts(attribute string) []string {
	cleanAttribute := partsRegex.ReplaceAllString(attribute, ".")
	cleanAttribute = strings.TrimRight(cleanAttribute, ".")
	return strings.Split(cleanAttribute, ".")
}
