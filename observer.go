package govuegui

import (
	"time"
)

// ObserverRefreshTime defines how often the routinge of the
// observer checks for changes.
var ObserverRefreshTime = 100 * time.Millisecond

type Observer struct {
	stop        chan int
	strings     map[*string]string
	subscribers []EventFunc
}

type EventFunc func(event, key, oldval string)

func NewObserver() *Observer {
	return &Observer{
		stop:    make(chan int),
		strings: make(map[*string]string),
	}
}

func (o *Observer) Start() {
	go func() {
		for {
			select {
			case <-time.Tick(time.Millisecond * 100):
				for p, val := range o.strings {
					if *p != val {
						o.Emmit("Value Changed", *p, val)
						o.strings[p] = *p
					}
				}
			case <-o.stop:
				return
			}
		}
	}()
}

func (o *Observer) Emmit(event, key, oldval string) {
	for _, f := range o.subscribers {
		f(event, key, oldval)
	}
}

func (o *Observer) Stop() {
	o.stop <- 0
}

func (o *Observer) Subscribe(ef EventFunc) {
	o.subscribers = append(o.subscribers, ef)
}

func (o *Observer) AddString(s *string) {
	o.strings[s] = *s

}
