package zero

import (
	"fmt"
	"time"
	"strconv"
	"strings"

	"github.com/tj/go-debug"
)

type numVal struct {
	Time time.Time
	Value int64
}

type numStore map[Key]numVal

type Num struct {
	store numStore
	debug debug.DebugFunction
}

func (n *Num) Add(key Key, value int64) {
	n.debug("adding key %v with value %d", key, value)
	n.store[key] = numVal{
		Time: time.Now(),
		Value: value,
	}
}

func (n *Num) Show(key Key) (int64, error) {
	n.debug("displaying value for key %v", key)
	if v, ok := n.store[key]; ok {
		n.debug("key found")
		return v.Value, nil
	}

	return 0, fmt.Errorf("key %s not found", key)
}

func (n *Num) All() [][]string {
	n.debug("displaying all values (%d)", len(n.store))
	vals := [][]string{}
	for k, v := range n.store {
		vals = append(vals, []string{
			string(k),
			strings.TrimSpace(strconv.FormatInt(v.Value, 10)),
			"number",
			v.Time.Format(time.RFC822),
		})
	}

	return vals
}

func (n *Num) Keys() []Key {
	n.debug("getting all keys for store (len: %d)", len(n.store))
	keys := make([]Key, 0, len(n.store))
	for k := range n.store {
		keys = append(keys, k)
	}

	return keys
}

func (n *Num) Count() int {
	n.debug("getting count of items in store")
	return len(n.store)
}

func (n *Num) Del(key Key) error {
	n.debug("removing value for key %v", key)
	if _, ok := n.store[key]; ok {
		n.debug("key found")
		delete(n.store, key)
		return nil
	}

	return fmt.Errorf("key %s not found", key)
}

func NewNum() *Num {
	return &Num{
		store: make(map[Key]numVal),
		debug: debug.Debug("mini:number"),
	}
}
