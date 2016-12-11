package zero

import (
	"fmt"
	"time"
	"strings"

	"github.com/tj/go-debug"
)

type strVal struct {
	Time  time.Time
	Value string
}

type strStore map[Key]strVal

type Str struct {
	store strStore
	debug debug.DebugFunction
}

func (s *Str) Add(key Key, value string) {
	s.debug("adding key %v with value %s", key, value)
	s.store[key] = strVal{
		Time: time.Now(),
		Value: value,
	}
}

func (s *Str) Show(key Key) (string, error) {
	s.debug("displaying value for key %v", key)
	if v, ok := s.store[key]; ok {
		s.debug("key found")
		return v.Value, nil
	}

	return "", fmt.Errorf("key %s not found", key)
}

func (s *Str) All() [][]string{
	s.debug("displaying all values (%d)", len(s.store))
	vals := [][]string{}
	for k, v := range s.store {
		r := []string{
			string(k),
			strings.TrimSpace(v.Value),
			"string",
			v.Time.Format(time.RFC822),
		}
		vals = append(vals, r)
	}
	return vals
}

func (s *Str) Keys() []Key {
	s.debug("getting all keys for store (len: %d)", len(s.store))
	keys := make([]Key, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}

	return keys
}

func (s *Str) Count() int {
	s.debug("getting count of items in store")
	return len(s.store)
}

func (s *Str) Del(key Key) error {
	s.debug("removing value for key %v", key)
	if _, ok := s.store[key]; ok {
		s.debug("key found")
		delete(s.store, key)
		return nil
	}

	return fmt.Errorf("key %s not found", key)
}

func NewStr() *Str {
	return &Str{
		store: make(strStore),
		debug: debug.Debug("zero:string"),
	}
}
