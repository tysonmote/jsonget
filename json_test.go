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
	WILDCARDS_JSON_PATH    = "test_json/wildcards.json"
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

func TestGetValues(t *testing.T) {
	testValues := map[string]string{
		// Given              Expected
		"foo":                "true",
		"bar":                "{\"baz\":5,\"biz\":5.5}",
		"bar.baz":            "5",
		"whoa[3]":            "spoon",
		"deep[1].peanuts[0]": "yum",
	}

	data, _ := LoadFile(GOOD_JSON_PATH)

	for attributeChain, expected := range testValues {
		values, err := data.GetValues(attributeChain)
		if err != nil {
			t.Fatal("Unexpected error: ", err, "while getting", attributeChain)
		}
		if values[0] != expected {
			t.Error("get didn't get the values for", attributeChain, "properly. Expected:", expected, "Got:", values[0])
		}
	}

	// Invalid access on non-objects
	_, err := data.GetValues("bar.baz.woo")
	expectedError := `Can't get "woo" on "bar.baz" because it is a float64`
	if err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
	_, err = data.GetValues("nope.nope")
	expectedError = `Cannot access "nope" on nil`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}

	// Invalid array access
	_, err = data.GetValues("deep[1].peanuts[10]")
	expectedError = `10 is outside the bounds of the 1 elements in "deep.1.peanuts"`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
	_, err = data.GetValues("deep[1].peanuts.yikes")
	expectedError = `"yikes" is not a valid index for the array at "deep.1.peanuts"`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
}

func TestGetValuesWildcards(t *testing.T) {
	testValues := map[string][]string{
		// Given              Expected
		"things.*":         []string{`{"names":["cool","sweet"],"size":2}`, `{"names":["rad"],"size":1}`, `{"names":["dude","bro","guys"],"size":3}`},
		"things.*.size":    []string{"2", "1", "3"},
		"things.*.names.0": []string{"cool", "rad", "dude"},
		"things.*.names.*": []string{"cool", "sweet", "rad", "dude", "bro", "guys"},
	}

	data, _ := LoadFile(WILDCARDS_JSON_PATH)

	for attributeChain, expected := range testValues {
		values, err := data.GetValues(attributeChain)
		if err != nil {
			t.Fatal("Unexpected error: ", err, "while getting", attributeChain)
		}
		if len(values) != len(expected) {
			t.Fatal("got", len(values), "values but expected", len(expected))
		}
		for i, value := range values {
			if value != expected[i] {
				t.Error("get didn't get the values for", attributeChain, "properly. Expected:", expected[i], "Got:", value)
			}
		}
	}

	// Wildcard on nil
	_, err := data.GetValues("foo.*")
	expectedError := `Cannot access "*" on nil`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}

	// Bad array access after wildcard
	_, err = data.GetValues("things.*.names[2]")
	expectedError = `2 is outside the bounds of the 2 elements in "things.*.names"`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}

	// Wildcard on non-array
	_, err = data.GetValues("things.*.size.*")
	expectedError = `Can't get "*" on "things.*.size" because it is a float64`
	if err == nil || err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
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
