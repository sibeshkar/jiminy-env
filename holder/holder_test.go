package holder

import (
	"testing"
	"time"
)

func TestValue(t *testing.T) {
	for _, h := range []Holder{
		NewMutexHolder(),
		NewChanHolder(),
	} {
		h.Set("foo")
		if got, want := h.Get(), "foo"; got != want {
			t.Errorf("Get() got %q, want %q", got, want)
		}

		h.Set("bar")
		if got, want := h.Get(), "bar"; got != want {
			t.Errorf("Get() got %q, want %q", got, want)
		}
		if got, want := h.Get(), "bar"; got != want {
			t.Errorf("Get() got %q, want %q", got, want)
		}

		// These tests are only supported by chanHolder.
		if ch, ok := h.(chanHolder); ok {
			// Attempt to Get while there is no value, should block indefinitely.
			ch.Set("") // clear the value.
			done := make(chan struct{})
			go func() {
				defer close(done)
				ch.Get()
			}()
			select {
			case <-time.After(100 * time.Millisecond): // expected.
			case <-done:
				t.Errorf("Get() returned, expected it to block")
			}

			// GetWithTimeout should time out eventually.
			if got, err := ch.GetWithTimeout(100 * time.Millisecond); err != ErrTimeout {
				t.Errorf("GetWithTimeout() got %q (err: %v), want timeout", got, err)
			}

			ch.Close()
		}
	}
}
