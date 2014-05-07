package model

import "testing"

type flCheck struct {
	ByteOffset int64
	Line       int64
	Char       int64
	Msg        string
}

func TestFindLine(t *testing.T) {
	file_string := []byte(`first line
second line


fifth line, after some blanks`)
	checks := []flCheck{
		flCheck{0, 1, 1, "Zero-offset has 1-indexed line/char"},
		flCheck{5, 1, 6, "Middle of first line"},
		flCheck{9, 1, 10, "End of first line"},
		flCheck{10, 1, 11, "Newline after line 1"},
		flCheck{11, 2, 1, "Start of line 2"},
		flCheck{15, 2, 5, "Middle of second line"},
		flCheck{21, 2, 11, "End of second line"},
		flCheck{22, 2, 12, "Newline after line 2"},
		flCheck{23, 3, 1, "Line 3"},
		flCheck{24, 4, 1, "Line 4"},
		flCheck{25, 5, 1, "Start of line 5"},
		flCheck{48, 5, 24, "The letter 'b'"},
		flCheck{53, 5, 29, "Last char ('s')"},
		flCheck{54, -1, -1, "After last char"},
		flCheck{154, -1, -1, "Long after last char"},
	}
    //t.Fatal("%v", file_string)
	for _, check := range checks {
		line, char := FindLine(file_string, check.ByteOffset)
		if line != check.Line || char != check.Char {
			t.Error(check)
			t.Fatalf("Got line %d, char %d", line, char)
		}
	}
}

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

// Failure cases
func TestGetJSONFromFile_NoSuchFile(t *testing.T) {
	_, err := GetJSONFromFile("nosuch.file")
	if err == nil {
		t.Fatal("GetJSONFromFile should fail on missing file")
	}
}

func TestGetJSONFromFile_BadJSON(t *testing.T) {
	_, err := GetJSONFromFile("broken.json")
	if err == nil {
		t.Fatal("GetJSONFromFile should fail on ill-formed JSON")
	}
    err_expected := `line 3, char 2: invalid character '}' looking for beginning of object key string`
    err_got := err.Error()
    if err_got != err_expected {
        t.Errorf("Expected: '%s'", err_expected)
        t.Errorf("Got: '%s'", err_got)
    }
}
