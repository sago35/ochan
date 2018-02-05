package ochan

type Ochan struct {
	In   []chan string
	out  chan string
	done chan struct{}
}

func NewOchan(ch chan string) *Ochan {
	o := &Ochan{
		out:  ch,
		done: make(chan struct{}),
	}

	go func(o *Ochan) {
		for _, ch := range o.In {
			for s := range ch {
				//fmt.Println("xxx", s)
				o.out <- s
			}
		}
		close(o.out)
		close(o.done)
	}(o)

	return o
}

func (o *Ochan) GetCh() chan<- string {
	ch := make(chan string, 100)
	o.In = append(o.In, ch)
	return ch
}

func (o *Ochan) Wait() error {
	<-o.done
	return nil
}
