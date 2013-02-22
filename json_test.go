package main

import (
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

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
}

//
// Tests
//

func TestJsonObjectFromFile(t *testing.T) {
	// Valid JSON
	data, err := JsonObjectFromFile(GOOD_JSON_PATH)
	assertNoError(t, err)
	if data["foo"] != true {
		t.Error("jsonFromFile didn't load valid JSON properly. Got:", data)
	}

	// Bad JSON
	data, err = JsonObjectFromFile(BAD_JSON_PATH)
	assertError(t, err)

	// Invalid path
	data, err = JsonObjectFromFile(NON_EXISTANT_JSON_PATH)
	assertError(t, err)
}

func TestValueToString(t *testing.T) {
	testValues := map[string]interface{}{
		// Expected            Given
		"Cool":                "Cool", // No surrounding quote marks
		"5":                   5,
		"1.23":                1.23,
		"[\"Cool\",\"Dude\"]": []string{"Cool", "Dude"},
		"{\"cool\":true}":     JsonObject{"cool": true},
	}

	for expected, given := range testValues {
		value, err := valueToString(given)
		assertNoError(t, err)
		if value != expected {
			t.Error("valueToString didn't convert a value properly. Expected:", expected, "Got:", value)
		}
	}

	defer reset()
	*printNulls = true
	value, err := valueToString(nil)
	assertNoError(t, err)
	if value != "null" {
		t.Error("When printNulls is true, valueToString should return \"null\" for nils.")
	}
}

func TestGetValue(t *testing.T) {
	testValues := map[string]string{
		// Expected                Given
		"true":                    "foo",
		"{\"baz\":5,\"biz\":5.5}": "bar",
		"5":                       "bar.baz",
		"":                        "nope.nope.nope",
	}

	data, _ := JsonObjectFromFile(GOOD_JSON_PATH)

	for expected, attributeChain := range testValues {
		value, err := data.GetValue(attributeChain)
		assertNoError(t, err)
		if value != expected {
			t.Error("get didn't get the values for", attributeChain, "properly. Expected:", expected, "Got:", value)
		}
	}

	// Invalid accessors
	_, err := data.GetValue("bar.baz.woo")
	expectedError := "Can't read woo attribute on bar.baz because it is not a JSON object."
	if err.Error() != expectedError {
		t.Error("Expected error message to be:", expectedError, "but got:", err)
	}
}

func TestGetValues(t *testing.T) {
	data, _ := JsonObjectFromFile(GOOD_JSON_PATH)

	attributes := []string{"foo", "bar.biz", "lol", "oh.no"}
	values, err := data.GetValues(attributes)
	assertNoError(t, err)
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
