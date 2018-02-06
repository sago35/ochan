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

	result := make(chan string, 100)
	o := NewOchan(result)
	n2 := runtime.NumGoroutine()
	if n1+1 != n2 {
		t.Errorf("NumGoroutine %d %d", n1, n2)
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
		t.Errorf("NumGoroutine %d %d", n2, n3)
	}

	o.Wait()
	close(result)

	i := 0
	for s := range result {
		if g, e := s, expected[i]; g != e {
			t.Errorf("got %q, want %q", g, e)
		}
		i++
	}

	n4 := runtime.NumGoroutine()
	if n1 != n4 {
		t.Errorf("NumGoroutine %d %d", n1, n4)
	}
}

func ExampleOchan() {
	result := make(chan string, 100)
	o := NewOchan(result)

	c1 := o.GetCh()
	c2 := o.GetCh()

	c1 <- "Hello c1"
	c2 <- "Hello c2"
	c1 <- "Bye c1"
	c2 <- "Bye c2"

	close(c1)
	close(c2)

	o.Wait()
	close(result)

	for s := range result {
		fmt.Println(s)
		// Output:
		// Hello c1
		// Bye c1
		// Hello c2
		// Bye c2
	}
}
