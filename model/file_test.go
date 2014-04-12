package model

import "testing"

// Success case
func TestGetJSONFromFile(t *testing.T) {
	value, err := GetJSONFromFile("demo.json")
	if err != nil {
		t.Fatal(err)
	}
	expected_keys := []string{"meta", "branches"}
	mv := value.(map[string]interface{})
	if len(mv) != len(expected_keys) {
		t.Fatalf(
			"Expected %d keys in data, got %d",
			len(expected_keys),
			len(mv),
		)
	}
	for _, key := range expected_keys {
		_, ok := mv[key]
		if !ok {
			t.Fatalf("Expected mv to contain key '%s'", key)
		}
	}
}
