package main

import (
	"fmt"
	"github.com/stvp/assert"
	"testing"
)

//
// Test helpers
//

const (
	GOOD_JSON_PATH         = "test_json/good.json"
	ARRAY_JSON_PATH        = "test_json/array.json"
	WILDCARDS_JSON_PATH    = "test_json/wildcards.json"
	BAD_JSON_PATH          = "test_json/bad.json"
	NON_EXISTANT_JSON_PATH = "test_json/lol.json"
)

func reset() {
	*printNulls = false
}

func testGetValues(t *testing.T, path string, expectedValues map[string][]string) {
	data, err := LoadFile(path)
	assert.Nil(t, err)

	for attributeChain, expected := range expectedValues {
		values, err := data.GetValues(attributeChain)
		assert.Nil(t, err)
		assert.Equal(t, len(expected), len(values))
		for i, value := range values {
			assert.Equal(t, expected[i], value)
		}
	}
}

func testGetValuesErrors(t *testing.T, path string, expectedErrors map[string]string) {
	data, err := LoadFile(path)
	assert.Nil(t, err)

	for attributeChain, expectedError := range expectedErrors {
		_, err = data.GetValues(attributeChain)
		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err.Error())
	}
}

//
// Tests
//

func TestJsonObjectFromFile(t *testing.T) {
	// Valid JSON
	data, err := LoadFile(GOOD_JSON_PATH)
	assert.Nil(t, err)
	assert.Equal(t, true, data.json.(map[string]interface{})["foo"])

	// Array JSON
	data, err = LoadFile(ARRAY_JSON_PATH)
	assert.Nil(t, err)
	assert.Equal(t, "Dude", (data.json.([]interface{})[0]).(map[string]interface{})["name"])

	// Bad JSON
	data, err = LoadFile(BAD_JSON_PATH)
	assert.NotNil(t, err)

	// Invalid path
	data, err = LoadFile(NON_EXISTANT_JSON_PATH)
	assert.NotNil(t, err)
}

func TestValueToString(t *testing.T) {
	testValues := map[string]interface{}{
		// Expected            Given
		"Cool":                "Cool", // No surrounding quote marks
		"5":                   5,
		"1.23":                1.23,
		"[\"Cool\",\"Dude\"]": []string{"Cool", "Dude"},
		"{\"cool\":true}":     map[string]interface{}{"cool": true},
	}

	for expected, given := range testValues {
		value, err := valueToString(given)
		assert.Nil(t, err)
		assert.Equal(t, expected, value)
	}

	defer reset()
	*printNulls = true
	value, err := valueToString(nil)
	assert.Nil(t, err)
	assert.Equal(t, "null", value, `when printNulls is true, nulls should be returned as "null"`)
}

func TestGetValues(t *testing.T) {
	testValues := map[string][]string{
		// Given              Expected
		"foo":                []string{"true"},
		"bar":                []string{"{\"baz\":5,\"biz\":5.5}"},
		"bar.baz":            []string{"5"},
		"whoa[3]":            []string{"spoon"},
		"deep[1].peanuts[0]": []string{"yum"},
	}
	testGetValues(t, GOOD_JSON_PATH, testValues)

	testErrors := map[string]string{
		// Given                 Expected error
		"bar.baz.woo":           `Can't get "woo" on "bar.baz" because it is a float64`,
		"nope.nope":             `Cannot access "nope" on nil`,
		"deep[1].peanuts[10]":   `10 is outside the bounds of the 1 elements in "deep.1.peanuts"`,
		"deep[1].peanuts.yikes": `"yikes" is not a valid index for the array at "deep.1.peanuts"`,
	}
	testGetValuesErrors(t, GOOD_JSON_PATH, testErrors)
}

func TestGetValuesWildcards(t *testing.T) {
	testValues := map[string][]string{
		// Given              Expected
		"things.*":         []string{`{"names":["cool","sweet"],"size":2}`, `{"names":["rad"],"size":1}`, `{"names":["dude","bro","guys"],"size":3}`},
		"things.*.size":    []string{"2", "1", "3"},
		"things.*.names.0": []string{"cool", "rad", "dude"},
		"things.*.names.*": []string{"cool", "sweet", "rad", "dude", "bro", "guys"},
	}
	testGetValues(t, WILDCARDS_JSON_PATH, testValues)

	testErrors := map[string]string{
		// Given             Expected error
		"foo.*":             `Cannot access "*" on nil`,
		"things.*.names[2]": `2 is outside the bounds of the 2 elements in "things.*.names"`,
		"things.*.size.*":   `Can't get "*" on "things.*.size" because it is a float64`,
	}
	testGetValuesErrors(t, WILDCARDS_JSON_PATH, testErrors)
}

func TestGetValuesArray(t *testing.T) {
	testValues := map[string][]string{
		// Given  Expected
		"*":      []string{`{"cool":false,"name":"Dude"}`, `{"cool":true,"name":"Sir"}`, `{"cool":false,"name":"Yo"}`},
		"*.name": []string{"Dude", "Sir", "Yo"},
	}
	testGetValues(t, ARRAY_JSON_PATH, testValues)

	testErrors := map[string]string{
		// Given  Expected error
		"foo.*": `Cannot access "foo" on the root array object`,
	}
	testGetValuesErrors(t, ARRAY_JSON_PATH, testErrors)
}

func TestSplitAttributeParts(t *testing.T) {
	parts := splitAttributeParts("foo[0][20].cool[2]")
	if len(parts) != 5 {
		t.Error("splitAttributeParts returned", len(parts), "parts, but we expected 5")
	}
	got := fmt.Sprintf("%v", parts)
	expected := "[foo 0 20 cool 2]"
	if got != expected {
		t.Error("splitAttributeParts returned", got, "but we expected", expected)
	}
}
