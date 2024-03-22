package cli

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache(time.Second)
	want := len(c.entries)
	if want != 0 {
		t.Fatalf("inital cache not empty")
	}

	c.Add("abc", []byte{})
	want = len(c.entries)
	if want == 0 {
		t.Fatal("cache did not add entry")
	}

	time.Sleep(time.Second * 2)
	want = len(c.entries)
	if want != 0 {
		t.Fatalf("cache not cleared")
	}
}

func TestCacheDupeAdd(t *testing.T) {
	c := NewCache(time.Second)

	c.Add("abc", []byte{})
	initial := len(c.entries)

	c.Add("abc", []byte{})
	post := len(c.entries)

	if initial != post {
		t.Fatal("duplicate add resulted in duplicate entries")
	}
}

func TestCacheAdd(t *testing.T) {
	c := NewCache(time.Second)

	initial := len(c.entries)
	c.Add("abc", []byte{})
	post := len(c.entries)

	if initial == post {
		t.Fatal("cache add failed to insert entry")
	}

	_, err := c.Get("abc")
	if err != nil {
		t.Fatal("cache does not have add entry")
	}
}

func TestCacheReap(t *testing.T) {
	c := NewCache(500 * time.Millisecond)

	initial := len(c.entries)
	c.Add("abc", []byte{})
	post := len(c.entries)

	if initial == post {
		t.Fatal("cache add failed to insert entry")
	}

	time.Sleep(time.Second)
	reap := len(c.entries)

	if reap == post {
		t.Fatal("cache failed to reap entry")
	}
}
