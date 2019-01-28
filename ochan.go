// Package ochan provides ordered chan.
package ochan

import (
	"sync"
)

// An Ochan is a structure for controlling the output order of channnels.
type Ochan struct {
	out  chan string
	in   chan chan string
	done chan struct{}
	wg   sync.WaitGroup
	size int
}

// Option ...
type Option func(*options)

type options struct {
	size int
}

var defaultOchanOption = options{
	size: 100,
}

// WithSize ...
func WithSize(s int) Option {
	return func(o *options) {
		o.size = s
	}
}

// NewOchan returns a new Ochan struct with specified buffer capacity.
func NewOchan(out chan string, opt ...Option) *Ochan {
	opts := defaultOchanOption
	for _, f := range opt {
		f(&opts)
	}

	o := &Ochan{
		out:  out,
		in:   make(chan chan string, opts.size),
		done: make(chan struct{}, 1),
		wg:   sync.WaitGroup{},
		size: opts.size,
	}

	go func(o *Ochan) {
		for {
			select {
			case ch, ok := <-o.in:
				if !ok {
					return
				}
				for s := range ch {
					o.out <- s
				}
				o.wg.Done()
			case <-o.done:
				return
			}
		}
	}(o)

	return o
}

// GetCh returns a next input channel. The input channel must be explicitly
// closed after use.
func (o *Ochan) GetCh() chan string {
	ch := make(chan string, o.size)
	o.in <- ch
	o.wg.Add(1)

	return ch
}

// SetSize sets the capacity of the channel returned by GetCh.
func (o *Ochan) SetSize(size int) {
	o.size = size
}

// Wait blocks until it retrieves data from all input channel. All input
// channels must be closed before calling this function.
func (o *Ochan) Wait() error {
	o.wg.Wait()
	return nil
}

// Close closed ochan's goroutine.
func (o *Ochan) Close() {
	close(o.done)
}
