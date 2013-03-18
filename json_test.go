package main

import (
	"fmt"
	"testing"
)

//
// Test helpers
//

const (
	GOOD_JSON_PATH         = "test_json/good.json"
	BAD_JSON_PATH          = "test_json/bad.json"
	NON_EXISTANT_JSON_PATH = "test_json/lol.json"
)

func reset() {
	*printNulls = false
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected error, but got nil instead.")
	}
}

//
// Tests
//

func TestJsonObjectFromFile(t *testing.T) {
	// Valid JSON
	data, err := LoadFile(GOOD_JSON_PATH)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if data.json.(map[string]interface{})["foo"] != true {
		t.Error("jsonFromFile didn't load valid JSON properly. Got:", data)
	}

	// Bad JSON
	data, err = LoadFile(BAD_JSON_PATH)
	assertError(t, err)

	// Invalid path
	data, err = LoadFile(NON_EXISTANT_JSON_PATH)
	assertError(t, err)
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
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		if value != expected {
			t.Error("valueToString didn't convert a value properly. Expected:", expected, "Got:", value)
		}
	}

	defer reset()
	*printNulls = true
	value, err := valueToString(nil)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if value != "null" {
		t.Error("When printNulls is true, valueToString should return \"null\" for nils.")
	}
}

func TestGetValue(t *testing.T) {
	testValues := map[string]string{
		// Expected               Given
		"true": "foo",
		"{\"baz\":5,\"biz\":5.5}": "bar",
		"5":                       "bar.baz",
		"":                        "nope.nope.nope",
		"spoon":                   "whoa[3]",
		"yum":                     "deep[1].peanuts[0]",
	}

	data, _ := LoadFile(GOOD_JSON_PATH)

	for expected, attributeChain := range testValues {
		value, err := data.GetValue(attributeChain)
		if err != nil {
			t.Fatal("Unexpected error: ", err, "while getting", attributeChain)
		}
		if value != expected {
			t.Error("get didn't get the values for", attributeChain, "properly. Expected:", expected, "Got:", value)
		}
	}

	// Invalid access on non-objects
	_, err := data.GetValue("bar.baz.woo")
	expectedError := "Can't get woo on bar.baz because it is a float64."
	if err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}

	// Invalid array access
	_, err = data.GetValue("deep[1].peanuts[10]")
	expectedError = "10 is outside the bounds of the 1 elements in deep.1.peanuts."
	if err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
	_, err = data.GetValue("deep[1].peanuts.yikes")
	expectedError = "deep.1.peanuts is an array, but yikes is not a valid index."
	if err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
}

func TestGetValues(t *testing.T) {
	data, _ := LoadFile(GOOD_JSON_PATH)

	attributes := []string{"foo", "bar.biz", "lol", "oh.no"}
	values, err := data.GetValues(attributes)
	if err != nil {
		t.Fatal("Unexpected error: ", err, "while getting values:", attributes)
	}
	if values[0] != "true" {
		t.Error("getValues returned:", values[0], "but we expected: true")
	}
	if values[1] != "5.5" {
		t.Error("getValues returned:", values[1], "but we expected: 5.5")
	}
	if values[2] != "" {
		t.Error("getValues returned:", values[2], "but we expected an empty string")
	}
	if values[3] != "" {
		t.Error("getValues returned:", values[3], "but we expected an empty string")
	}
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
