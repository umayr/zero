package zero

import (
	"fmt"

	"github.com/tj/go-debug"
)

type Key string

type Args struct {
	Key   Key
	Type  string
	Value interface{}
}

const (
	String = "string"
	Number = "number"
	Array  = "array"
)

const (
	RowSplitter    = "|"
	ColumnSplitter = ","
)

var d = debug.Debug("zero")

var index = make(map[Key]string)

func add(key Key, kind string) {
	d("associating key %s to type %s", key, kind)
	index[key] = kind
}

func which(key Key) (string, error) {
	d("pulling type of key %s", key)
	if v, exists := index[key]; exists {
		return v, nil
	}
	return "", fmt.Errorf("key %s not found", key)
}
