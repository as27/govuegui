package govuegui

import (
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
	stop        chan int
	strings     []*stringValue
	subscribers []EventFunc
	RefreshTime time.Duration
}

// EventFunc is used to define a function which is called
// when a value changes
type EventFunc func(event, key, oldval string)

// NewObserver creates a empty Observer
func NewObserver() *Observer {
	return &Observer{
		RefreshTime: ObserverRefreshTime,
		stop:        make(chan int),
	}
}

// Start to begin watching the varaibles
func (o *Observer) Start() {
	go func() {
		for {
			select {
			case <-time.Tick(time.Millisecond * 100):
				for _, sval := range o.strings {
					if *sval.pointer != sval.value {
						o.emmit(ObserverValueChanged, sval.key, sval.value)
						sval.value = *sval.pointer
					}
				}
			case <-o.stop:
				return
			}
		}
	}()
}

func (o *Observer) emmit(event, key, oldval string) {
	for _, f := range o.subscribers {
		f(event, key, oldval)
	}
}

// Stop the observer watching the variables
func (o *Observer) Stop() {
	o.stop <- 0
}

// Subscribe a func type EventFunc to the observer. All subscribed
// functions are called, when a event happens
func (o *Observer) Subscribe(ef EventFunc) {
	o.subscribers = append(o.subscribers, ef)
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
