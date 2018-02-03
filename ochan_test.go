package ochan

import (
	"fmt"
	"runtime"
	"testing"
)

func TestOchanBasic(t *testing.T) {
	n1 := runtime.NumGoroutine()

	expected := []string{
		"c1-1",
		"c1-2",
		"c1-3",
		"c2-1",
		"c2-2",
		"c3-1",
		"c3-2",
	}

	o := NewOchan()
	n2 := runtime.NumGoroutine()
	if n1+1 != n2 {
		t.Errorf("NumGoroutine")
	}

	c1 := o.GetCh()
	c1 <- "c1-1"
	c1 <- "c1-2"

	c2 := o.GetCh()

	c3 := o.GetCh()
	c3 <- "c3-1"
	c2 <- "c2-1"
	c3 <- "c3-2"
	c2 <- "c2-2"
	c1 <- "c1-3"

	close(c1)
	close(c2)
	close(c3)

	n3 := runtime.NumGoroutine()
	if n2 != n3 {
		t.Errorf("NumGoroutine")
	}

	i := 0
	for s := range o.Out {
		if g, e := s, expected[i]; g != e {
			t.Errorf("got %q, want %q", g, e)
		}
		i++
	}
	o.Wait()

	n4 := runtime.NumGoroutine()
	if n1 != n4 {
		t.Errorf("NumGoroutine", n1, n4)
	}
}

func ExampleOchan() {
	o := NewOchan()

	c1 := o.GetCh()
	c2 := o.GetCh()

	c1 <- "Hello c1"
	c2 <- "Hello c2"
	c1 <- "Bye c1"
	c2 <- "Bye c2"

	close(c1)
	close(c2)

	for s := range o.Out {
		fmt.Println(s)
		// Output:
		// Hello c1
		// Bye c1
		// Hello c2
		// Bye c2
	}
	o.Wait()
}
