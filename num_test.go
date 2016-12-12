package zero

import (
	"strings"
	"testing"
	"time"
)

func TestNewNum(t *testing.T) {
	num := NewNum()
	if num.store == nil {
		t.Error("store is not initialised")
	}
}

func TestNumAdd(t *testing.T) {
	num := NewNum()

	num.Add(Key("foo"), 12)
	num.Add(Key("bar"), 34)

	if len(num.store) != 2 {
		t.Error("values are not being added to the store")
	}

	val := num.store[Key("foo")]
	if val.Time.IsZero() {
		t.Error("time is not being set on entry")
	}

	if val.Value == 0 {
		t.Error("value is not being set on entry")
	}
}

func TestNumShow(t *testing.T) {
	num := NewNum()

	num.store[Key("foo")] = numVal{time.Now(), 12}

	s, _ := num.Show(Key("foo"))
	if s != 12 {
		t.Error("returning invalid value")
	}

	_, err := num.Show(Key("bar"))
	if err == nil {
		t.Error("no error when key is not in the store")
	}
}

func TestNumAll(t *testing.T) {
	num := NewNum()
	num.store[Key("foo")] = numVal{time.Now(), 12}
	num.store[Key("bar")] = numVal{time.Now(), 34}

	all := num.All()
	if len(all) != 2 {
		t.Error("length should be equal to two")
	}

	for _, v := range all {
		if v[0] == "" {
			t.Error("key is not being returned")
		}

		if v[1] == "" {
			t.Error("value is not being returned")
		}

		if v[2] != "number" {
			t.Error("type is not number")
		}

		if v[3] == "" {
			t.Error("time is not being returned")
		}
	}
}

func TestNumKeys(t *testing.T) {
	num := NewNum()
	num.store[Key("foo")] = numVal{time.Now(), 12}
	num.store[Key("bar")] = numVal{time.Now(), 34}

	keys := num.Keys()
	if len(keys) != 2 {
		t.Error("length should be equal to two")
	}

	joined := func() string {
		v := make([]string, len(keys))
		for _, k := range keys {
			v = append(v, string(k))
		}
		return strings.Join(v, " ")
	}()

	if !strings.Contains(joined, "foo") {
		t.Error("keys should contain foo")
	}

	if !strings.Contains(joined, "bar") {
		t.Error("keys should contain bar")
	}
}

func TestNumCount(t *testing.T) {
	num := NewNum()
	num.store[Key("foo")] = numVal{time.Now(), 12}
	num.store[Key("bar")] = numVal{time.Now(), 34}

	count := num.Count()
	if count != 2 {
		t.Error("count should be equal to two")
	}
}

func TestNumDel(t *testing.T) {
	num := NewNum()
	num.store[Key("foo")] = numVal{time.Now(), 12}
	num.store[Key("bar")] = numVal{time.Now(), 34}

	num.Del(Key("foo"))
	if len(num.store) != 1 {
		t.Error("length should be equal to one")
	}

	num.Del(Key("bar"))
	if len(num.store) != 0 {
		t.Error("length should be equal to zero")
	}

	err := num.Del(Key("baz"))
	if err == nil {
		t.Error("no error when key is not found")
	}
}
