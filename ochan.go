package ochan

type Ochan struct {
	out  chan string
	in   chan chan string
	done chan struct{}
}

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

func (o *Ochan) GetCh() chan string {
	ch := make(chan string, 100)
	o.in <- ch

	return ch
}

func (o *Ochan) Wait() error {
	close(o.in)
	<-o.done
	return nil
}
