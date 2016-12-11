package zero

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tj/go-debug"
)

type arrVal struct {
	Time  time.Time
	Items []string

	strItems map[string]string
	numItems map[string]int64
}

type arrStore map[Key]arrVal

type Arr struct {
	store arrStore
	debug debug.DebugFunction
}

func (a *Arr) Add(key Key, values []interface{}) error {
	v := arrVal{Time: time.Now(), strItems: map[string]string{}, numItems: map[string]int64{}}

	for i, value := range values {
		switch typeOf(value) {
		case String:
			s, ok := value.(string)
			if !ok {
				return fmt.Errorf("unable to parse value")
			}

			k := fmt.Sprintf("%s.#%d", key, i)
			v.strItems[k] = s
			v.Items = append(v.Items, k)
			break
		case Number:
			n, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 0)
			if err != nil {
				return err
			}

			k := fmt.Sprintf("%s.#%d", key, i)
			v.numItems[k] = n
			v.Items = append(v.Items, k)

			break
		}
	}

	a.debug("adding key %v with value %v", key, values)
	a.store[key] = v
	return nil
}

func (a *Arr) Has(key Key) bool {
	for k, v := range a.store {
		a.debug("key: %s value: %v", k, v)
	}
	_, exists := a.store[key]
	return exists
}

func (a *Arr) Push(key Key, value interface{}) (int, error) {
	if !a.Has(key) {
		return 0, fmt.Errorf("key %s not found", key)
	}

	v := a.store[key]
	l := 0
	switch typeOf(value) {
	case String:
		s, ok := value.(string)
		if !ok {
			return l, fmt.Errorf("unable to parse value")
		}

		l = len(v.strItems) + len(v.numItems)
		k := fmt.Sprintf("%s.#%d", key, l)
		v.strItems[k] = s
		v.Items = append(v.Items, k)

		break
	case Number:
		n, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 0)
		if err != nil {
			return l, err
		}

		l = len(v.strItems) + len(v.numItems)
		k := fmt.Sprintf("%s.#%d", key, l)
		v.numItems[k] = n
		v.Items = append(v.Items, k)

		break
	}
	a.store[key] = v
	return l + 1, nil
}

func (a *Arr) Pop(key Key) (interface{}, error) {
	if !a.Has(key) {
		return nil, fmt.Errorf("key %s not found", key)
	}

	v := a.store[key]
	if len(v.Items) == 0 {
		return nil, fmt.Errorf("array is already empty")
	}

	var (
		last string
		ret interface{}
	)

	v.Items, last = v.Items[:len(v.Items)-1], v.Items[len(v.Items)-1]
	if val, has := v.strItems[last]; has {
		ret = val
		delete(v.strItems, last)
	} else if val, has := v.numItems[last]; has {
		ret = val
		delete(v.numItems, last)
	}

	a.store[key] = v
	return ret, nil
}

func (a *Arr) Show(key Key) ([]interface{}, error) {
	if !a.Has(key) {
		return nil, fmt.Errorf("key %s not found", key)
	}
	a.debug("displaying value for key %v", key)

	v := a.store[key]
	values := []interface{}{}

	for _, k := range v.Items {
		if str, has := v.strItems[k]; has {
			values = append(values, str)
		} else if num, has := v.numItems[k]; has {
			values = append(values, num)
		}
	}

	return values, nil
}

func (a *Arr) All() [][]string {
	a.debug("displaying all values (%d)", len(a.store))
	vals := [][]string{}
	for k, v := range a.store {
		val, _ := a.Show(k)
		r := []string{
			string(k),
			strings.TrimSpace(fmt.Sprintf("%v", val)),
			"array",
			v.Time.Format(time.RFC822),
		}

		vals = append(vals, r)
	}
	return vals
}

func (a *Arr) Keys() []Key {
	a.debug("getting all keys for store (len: %d)", len(a.store))
	keys := make([]Key, 0, len(a.store))
	for k := range a.store {
		keys = append(keys, k)
	}

	return keys
}

func (a *Arr) Count() int {
	a.debug("getting count of items in store")
	return len(a.store)
}

func (a *Arr) Del(key Key) error {
	a.debug("removing value for key %v", key)
	if _, ok := a.store[key]; ok {
		a.debug("key found")
		delete(a.store, key)
		return nil
	}

	return fmt.Errorf("key %s not found", key)
}

func typeOf(v interface{}) string {
	if s, ok := v.(string); ok {
		if _, err := strconv.ParseInt(s, 10, 0); err == nil {
			return Number
		}

		return String
	}
	return ""
}

func NewArr() *Arr {
	return &Arr{
		store: make(arrStore),
		debug: debug.Debug("zero:array"),
	}
}
