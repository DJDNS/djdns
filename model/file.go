package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func FindLine(text []byte, offset int64) (line int64, char int64) {
	// TODO: more line ending support using bytes.Replace
    var pos int64
    line = 1
    char = 1
    if (offset >= int64(len(text))) {
        line, char = -1, -1
        return
    }
    for pos < offset {
        if text[pos] == '\n' {
            line++
            char = 1
        } else {
            char++
        }

        pos++
    }
	return
}

func GetJSONFromFile(filename string) (interface{}, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			line, col := FindLine(raw, serr.Offset)
			err = fmt.Errorf("line %d, char %d: %s", line, col, err)
		}
	}
	return result, err
}
