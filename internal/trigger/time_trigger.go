package trigger

import (
	"sync"
	"time"
)

// TimeTrigger implements the Trigger interface using a timer
type TimeTrigger struct {
	interval time.Duration
	ticker   *time.Ticker
	done     chan bool
	handler  Handler
	running  bool
	mutex    sync.Mutex
}

// NewTimeTrigger creates a new TimeTrigger with the specified interval
func NewTimeTrigger(interval time.Duration) *TimeTrigger {
	return &TimeTrigger{
		interval: interval,
		done:     make(chan bool),
	}
}

// Start starts the trigger with the specified handler
func (t *TimeTrigger) Start(handler Handler) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.running {
		return nil
	}

	t.handler = handler
	t.ticker = time.NewTicker(t.interval)
	t.running = true

	go func() {
		// Execute the handler immediately on start
		if err := t.handler(); err != nil {
		}

		for {
			select {
			case <-t.ticker.C:
				if err := t.handler(); err != nil {
				}
			case <-t.done:
				return
			}
		}
	}()

	return nil
}

// Stop stops the trigger
func (t *TimeTrigger) Stop() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.running {
		return nil
	}

	t.ticker.Stop()
	t.done <- true
	t.running = false

	return nil
}
