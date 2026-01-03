package singleflight

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestWrapper_Do(t *testing.T) {
	synctest.Test(t, func(*testing.T) {
		g := New(context.Background(), nil)

		var calls int
		fn := func() (any, error) {
			calls++
			time.Sleep(10 * time.Millisecond)
			return "bar", nil
		}

		done := make(chan struct{}, 10)
		for i := 0; i < 10; i++ {
			go func() {
				v, err, _ := g.Do("key", fn)
				if err != nil {
					t.Errorf("Do error: %v", err)
				}
				if v != "bar" {
					t.Errorf("got %v, want %v", v, "bar")
				}
				done <- struct{}{}
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		if calls != 1 {
			t.Errorf("number of calls = %d, want 1", calls)
		}
	})
}
