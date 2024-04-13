package cache

import "testing"

func TestMapCache(t *testing.T) {
	mc := NewMapCache()
	mc.Set("key", "value")
	if v, err := mc.Get("key"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(v)
	}
}
