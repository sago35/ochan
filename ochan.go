// Package ochan provides ordered chan.
package ochan

import (
	"sync"
)

// An Ochan is a structure for controlling the output order of channels.
type Ochan[T any] struct {
	out  chan T
	in   chan chan T
	done chan struct{}
	wg   sync.WaitGroup
	size int
}

// NewOchan returns a new Ochan struct with specified buffer capacity.
func NewOchan[T any](out chan T, size int) *Ochan[T] {
	o := &Ochan[T]{
		out:  out,
		in:   make(chan chan T, size),
		done: make(chan struct{}, 1),
		wg:   sync.WaitGroup{},
		size: size,
	}

	go func(o *Ochan[T]) {
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
func (o *Ochan[T]) GetCh() chan T {
	ch := make(chan T, o.size)
	o.in <- ch
	o.wg.Add(1)

	return ch
}

// SetSize sets the capacity of the channel returned by GetCh.
func (o *Ochan[T]) SetSize(size int) {
	o.size = size
}

// Wait blocks until it retrieves data from all input channel. All input
// channels must be closed before calling this function.
func (o *Ochan[T]) Wait() error {
	o.wg.Wait()
	return nil
}

// Close closed ochan's goroutine.
func (o *Ochan[T]) Close() {
	close(o.done)
}
