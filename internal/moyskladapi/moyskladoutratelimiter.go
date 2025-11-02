package moyskladapi

import "time"

type Ratelimiter struct {
	ticker *time.Ticker
	ch     chan struct{}
}

func (r *Ratelimiter) run() {
	for range r.ticker.C {
		select {
		case r.ch <- struct{}{}:
		default:
		}
	}
}
func NewRatelimiter(limit int, interval time.Duration) *Ratelimiter {
	r := &Ratelimiter{
		ticker: time.NewTicker(interval / time.Duration(limit)),
		ch:     make(chan struct{}, limit),
	}
	go r.run()
	return r
}

func (r *Ratelimiter) Wait() {
	<-r.ch
}

func (r *Ratelimiter) Stop() {
	r.ticker.Stop()
	close(r.ch)
}

func (r *Ratelimiter) Chan() <-chan struct{} {
	return r.ch
}
