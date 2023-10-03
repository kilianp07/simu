package utils

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(path string, target interface{}) (err error) {
	var (
		content []byte
	)

	content, err = os.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, target)
	return
}
