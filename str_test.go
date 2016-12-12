package zero

import (
	"strings"
	"testing"
	"time"
)

func TestNewStr(t *testing.T) {
	str := NewStr()
	if str.store == nil {
		t.Error("store is not initialised")
	}
}

func TestStrAdd(t *testing.T) {
	str := NewStr()

	str.Add(Key("foo"), "bar")
	str.Add(Key("bar"), "baz")

	if len(str.store) != 2 {
		t.Error("values are not being added to the store")
	}

	val := str.store[Key("foo")]
	if val.Time.IsZero() {
		t.Error("time is not being set on entry")
	}

	if val.Value == "" {
		t.Error("value is not being set on entry")
	}
}

func TestStrShow(t *testing.T) {
	str := NewStr()

	str.store[Key("foo")] = strVal{time.Now(), "bar"}

	s, _ := str.Show(Key("foo"))
	if s != "bar" {
		t.Error("returning invalid value")
	}

	_, err := str.Show(Key("bar"))
	if err == nil {
		t.Error("no error when key is not in the store")
	}
}

func TestStrAll(t *testing.T) {
	str := NewStr()
	str.store[Key("foo")] = strVal{time.Now(), "bar"}
	str.store[Key("bar")] = strVal{time.Now(), "baz"}

	all := str.All()
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

		if v[2] != "string" {
			t.Error("type is not string")
		}

		if v[3] == "" {
			t.Error("time is not being returned")
		}
	}
}

func TestStrKeys(t *testing.T) {
	str := NewStr()
	str.store[Key("foo")] = strVal{time.Now(), "bar"}
	str.store[Key("bar")] = strVal{time.Now(), "baz"}

	keys := str.Keys()
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

func TestStrCount(t *testing.T) {
	str := NewStr()
	str.store[Key("foo")] = strVal{time.Now(), "bar"}
	str.store[Key("bar")] = strVal{time.Now(), "baz"}

	count := str.Count()
	if count != 2 {
		t.Error("count should be equal to two")
	}
}

func TestStrDel(t *testing.T) {
	str := NewStr()
	str.store[Key("foo")] = strVal{time.Now(), "bar"}
	str.store[Key("bar")] = strVal{time.Now(), "baz"}

	str.Del(Key("foo"))
	if len(str.store) != 1 {
		t.Error("length should be equal to one")
	}

	str.Del(Key("bar"))
	if len(str.store) != 0 {
		t.Error("length should be equal to zero")
	}

	err := str.Del(Key("baz"))
	if err == nil {
		t.Error("no error when key is not found")
	}
}
