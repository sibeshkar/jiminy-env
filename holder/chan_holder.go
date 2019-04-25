package holder

import (
	"errors"
	"time"
)

type chanHolder struct {
	setValCh chan string
	getValCh chan string
	closeCh  chan struct{}
}

// NewChanHolder returns a new Holder backed by Channels.
func NewChanHolder() Holder {
	h := chanHolder{
		setValCh: make(chan string),
		getValCh: make(chan string),
		closeCh:  make(chan struct{}),
	}
	go h.mux()
	return h
}

func (h chanHolder) mux() {
	var value string
	for {
		// if the value is empty, only accept setting, or closing.
		if value == "" {
			select {
			case <-h.closeCh: // we also need to handle closing here!
				close(h.setValCh)
				close(h.getValCh)
				return
			case value = <-h.setValCh:
				continue
			}
		}
		// once the value is non-empty, it can be set or gotten
		// as normal.
		select {
		case value = <-h.setValCh: // set the current value.
		case h.getValCh <- value: // send the current value.
		case <-h.closeCh: // closing, time to clean up!
			close(h.setValCh)
			close(h.getValCh)
			return
		}
	}
}

func (h chanHolder) Get() string {
	return <-h.getValCh
}

func (h chanHolder) Set(s string) {
	h.setValCh <- s
}

// Close closes the Holder, making calls to Set panic, and calls to Get return
// the last-set value.
func (h chanHolder) Close() {
	close(h.closeCh)
}

// ErrTimeout is the error returned by GetWithTimeout if the value
// was not provided before the given timeout.
var ErrTimeout = errors.New("timeout waiting for value")

// GetWithTimeout attempts to get the value, or returns ErrTimeout
// if getting it takes too long.
func (h chanHolder) GetWithTimeout(d time.Duration) (string, error) {
	select {
	case <-time.After(d):
		return "", ErrTimeout
	case v := <-h.getValCh:
		return v, nil
	}
}
