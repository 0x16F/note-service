package jsonmap

import "github.com/goccy/go-json"

func (m Map) String() string {
	encoded, _ := json.Marshal(&m)
	return string(encoded)
}
