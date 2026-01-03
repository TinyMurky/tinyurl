package singleflight

import (
	"testing"

	"github.com/TinyMurky/tinyurl/pkg/singleflight"
)

func TestGroup_Do_Prefix(t *testing.T) {
	mockGroup := &singleflight.MockGroup{}
	g := New(mockGroup)

	id := "test-id"
	expectedKey := "shorturl:test-id"

	called := false
	mockGroup.DoFunc = func(key string, fn func() (any, error)) (v any, err error, shared bool) {
		called = true
		if key != expectedKey {
			t.Errorf("Do called with key %q, want %q", key, expectedKey)
		}
		v, err = fn()
		return v, err, false
	}

	g.Do(id, func() (any, error) {
		return "result", nil
	})

	if !called {
		t.Error("underlying Do was not called")
	}
}
