package zero

import (
	"strings"
	"testing"
	"time"
)

func TestNewArr(t *testing.T) {
	arr := NewArr()
	if arr.store == nil {
		t.Error("store is not initialised")
	}
}

func TestArrAdd(t *testing.T) {
	arr := NewArr()

	err0 := arr.Add(Key("foo"), []interface{}{"1000", "2000"})
	err1 := arr.Add(Key("bar"), []interface{}{"foo", "bar"})
	err2 := arr.Add(Key("baz"), []interface{}{"500", "baz"})

	if err0 != nil || err1 != nil || err2 != nil {
		t.Error("error while adding items")
	}

	if len(arr.store) != 3 {
		t.Error("values are not being added to the store")
	}

	v0 := arr.store[Key("foo")]
	if v0.Time.IsZero() {
		t.Error("time is not being set on entry")
	}

	if len(v0.Items) != 2 {
		t.Error("value is not being set on entry")
	}

	if len(v0.numItems) != 2 {
		t.Error("numeric value is not being set on entry")
	}

	if len(v0.strItems) != 0 {
		t.Error("numeric value is not being set on entry")
	}

	v1 := arr.store[Key("bar")]
	if v1.Time.IsZero() {
		t.Error("time is not being set on entry")
	}

	if len(v1.Items) != 2 {
		t.Error("value is not being set on entry")
	}

	if len(v1.numItems) != 0 {
		t.Error("string value is not being set on entry")
	}

	if len(v1.strItems) != 2 {
		t.Error("string value is not being set on entry")
	}

	v2 := arr.store[Key("baz")]
	if v2.Time.IsZero() {
		t.Error("time is not being set on entry")
	}

	if len(v2.Items) != 2 {
		t.Error("value is not being set on entry")
	}

	if len(v2.numItems) != 1 {
		t.Error("mix value is not being set on entry")
	}

	if len(v2.strItems) != 1 {
		t.Error("mix value is not being set on entry")
	}
}

func TestArrHas(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand.#0", "rand.#1"},
		strItems: map[string]string{
			"rand.#0": "hello",
		},
		numItems: map[string]int64{
			"rand.#1": 20394,
		},
	}

	r0 := arr.Has(Key("foo"))
	if !r0 {
		t.Error("store should have element with foo key")
	}

	r1 := arr.Has(Key("bar"))
	if r1 {
		t.Error("store should not have element with bar key")
	}
}

func TestArrShow(t *testing.T) {
	arr := NewArr()

	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand.#0", "rand.#1"},
		strItems: map[string]string{
			"rand.#0": "hello",
		},
		numItems: map[string]int64{
			"rand.#1": 20394,
		},
	}

	v, _ := arr.Show(Key("foo"))
	if len(v) != 2 {
		t.Error("length should be equal to two")
	}

	_, err := arr.Show(Key("bar"))
	if err == nil {
		t.Error("no error when key is not in the store")
	}
}

func TestArrPush(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"foo.#0", "foo.#1"},
		strItems: map[string]string{
			"foo.#0": "hello",
		},
		numItems: map[string]int64{
			"foo.#1": 20394,
		},
	}

	old := len(arr.store[Key("foo")].Items)
	new, _ := arr.Push(Key("foo"), "pushed value")

	if old >= new {
		t.Error("new length should be greater than old one")
	}

	if arr.store[Key("foo")].strItems["foo.#2"] != "pushed value" {
		t.Error("pushed value should be equal to given value")
	}
}

func TestArrPop(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"foo.#0", "foo.#1"},
		strItems: map[string]string{
			"foo.#0": "hello",
		},
		numItems: map[string]int64{
			"foo.#1": 20394,
		},
	}

	arr.Pop(Key("foo"))
	if len(arr.store[Key("foo")].Items) != 1 {
		t.Error("length should be equal to one")
	}

	arr.Pop(Key("foo"))
	if len(arr.store[Key("foo")].Items) != 0 {
		t.Error("length should be equal to zero")
	}

	_, err := arr.Pop(Key("foo"))
	if err == nil {
		t.Error("no error when popping item from empty array")
	}
}

func TestArrAll(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand0.#0", "rand0.#1"},
		strItems: map[string]string{
			"rand0.#0": "hello",
		},
		numItems: map[string]int64{
			"rand0.#1": 20394,
		},
	}
	arr.store[Key("bar")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand1.#0", "rand1.#1"},
		strItems: map[string]string{
			"rand1.#0": "world",
		},
		numItems: map[string]int64{
			"rand1.#1": 49302,
		},
	}

	all := arr.All()
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

		if v[2] != "array" {
			t.Error("type is not array")
		}

		if v[3] == "" {
			t.Error("time is not being returned")
		}
	}
}

func TestArrKeys(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand0.#0", "rand0.#1"},
		strItems: map[string]string{
			"rand0.#0": "hello",
		},
		numItems: map[string]int64{
			"rand0.#1": 20394,
		},
	}
	arr.store[Key("bar")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand1.#0", "rand1.#1"},
		strItems: map[string]string{
			"rand1.#0": "world",
		},
		numItems: map[string]int64{
			"rand1.#1": 49302,
		},
	}

	keys := arr.Keys()
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

func TestArrCount(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand0.#0", "rand0.#1"},
		strItems: map[string]string{
			"rand0.#0": "hello",
		},
		numItems: map[string]int64{
			"rand0.#1": 20394,
		},
	}
	arr.store[Key("bar")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand1.#0", "rand1.#1"},
		strItems: map[string]string{
			"rand1.#0": "world",
		},
		numItems: map[string]int64{
			"rand1.#1": 49302,
		},
	}

	count := arr.Count()
	if count != 2 {
		t.Error("count should be equal to two")
	}
}

func TestArrDel(t *testing.T) {
	arr := NewArr()
	arr.store[Key("foo")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand0.#0", "rand0.#1"},
		strItems: map[string]string{
			"rand0.#0": "hello",
		},
		numItems: map[string]int64{
			"rand0.#1": 20394,
		},
	}
	arr.store[Key("bar")] = arrVal{
		Time:  time.Now(),
		Items: []string{"rand1.#0", "rand1.#1"},
		strItems: map[string]string{
			"rand1.#0": "world",
		},
		numItems: map[string]int64{
			"rand1.#1": 49302,
		},
	}

	arr.Del(Key("foo"))
	if len(arr.store) != 1 {
		t.Error("length should be equal to one")
	}

	arr.Del(Key("bar"))
	if len(arr.store) != 0 {
		t.Error("length should be equal to zero")
	}

	err := arr.Del(Key("baz"))
	if err == nil {
		t.Error("no error when key is not found")
	}
}
