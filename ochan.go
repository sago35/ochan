// Package ochan provides ordered chan.
package ochan

// An Ochan is a structure for controlling the output order of channnels.
type Ochan struct {
	out  chan string
	in   chan chan string
	done chan struct{}
}

// NewOchan returns a new Ochan struct.
func NewOchan(out chan string) *Ochan {
	o := &Ochan{
		out:  out,
		in:   make(chan chan string, 100),
		done: make(chan struct{}),
	}

	go func(o *Ochan) {
		for {
			select {
			case ch, ok := <-o.in:
				if !ok {
					close(o.done)
					return
				}
				for s := range ch {
					o.out <- s
				}
			}
		}
	}(o)

	return o
}

// GetCh returns a next input channel. The input channel must be explicitly
// closed after use.
func (o *Ochan) GetCh() chan string {
	ch := make(chan string, 100)
	o.in <- ch

	return ch
}

// Wait blocks until it retrieves data from all input channel. All input
// channels must be closed before calling this function.
func (o *Ochan) Wait() error {
	close(o.in)
	<-o.done
	return nil
}
