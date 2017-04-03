package observer

import (
	"sync"
	"time"
)

const (
	ObserverValueChanged = "Value Changed"
)

// ObserverRefreshTime defines how often the routinge of the
// observer checks for changes.
var ObserverRefreshTime = 100 * time.Millisecond

type stringValue struct {
	key     string
	value   string
	pointer *string
}

// Observer watches variables if they change
type Observer struct {
	runs        bool
	stop        chan int
	mutex       sync.Mutex
	strings     []*stringValue
	subscribers []chan stringValue
	RefreshTime time.Duration
}

// EventFunc is used to define a function which is called
// when a value changes
type EventFunc func(event string, sv chan stringValue)

// NewObserver creates a empty Observer
func NewObserver() *Observer {
	return &Observer{
		RefreshTime: ObserverRefreshTime,
		stop:        make(chan int),
	}
}

// Start to begin watching the varaibles
func (o *Observer) Start() {
	if !o.runs {
		go func() {
			for {
				select {
				case <-time.Tick(time.Millisecond * 100):
					for _, sval := range o.strings {
						if *sval.pointer != sval.value {
							o.mutex.Lock()
							o.emmit(ObserverValueChanged, *sval)

							sval.value = *sval.pointer
							o.mutex.Unlock()
						}
					}
				case <-o.stop:
					return
				}
			}
		}()
	}
}

func (o *Observer) emmit(event string, sv stringValue) {
	for _, f := range o.subscribers {
		f <- sv
	}
}

// Stop the observer watching the variables
func (o *Observer) Stop() {
	o.stop <- 0
}

// Subscribe a func type EventFunc to the observer. All subscribed
// functions are called, when a event happens
func (o *Observer) Subscribe(svc chan stringValue) {
	o.subscribers = append(o.subscribers, svc)
}

// AddString to the observer. The key is used to identify that
// variable, when a event emitts data
func (o *Observer) AddString(key string, s *string) {
	o.strings = append(
		o.strings,
		&stringValue{
			key:     key,
			value:   *s,
			pointer: s,
		},
	)

}
