package main

import (
	"encoding/json"
	"io/ioutil"
)

func GetJSONFromFile(filename string) (interface{}, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(raw, &result)
	return result, err
}
