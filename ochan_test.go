package ochan

import (
	"fmt"
	"runtime"
	"testing"
	"time"
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

	expected2 := []string{
		"c4-1",
		"c4-2",
		"c5-1",
		"c5-2",
	}

	result := make(chan string, 100)
	o := NewOchan(result, 100)
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

	for i := 0; i < len(expected); i++ {
		s := <-result
		if g, e := s, expected[i]; g != e {
			t.Errorf("%d got %q, want %q", i, g, e)
		}
	}

	n4 := runtime.NumGoroutine()
	if n1+1 != n4 {
		t.Errorf("NumGoroutine %d %d", n1+1, n4)
	}

	c4 := o.GetCh()
	c5 := o.GetCh()

	c4 <- "c4-1"
	c5 <- "c5-1"
	c4 <- "c4-2"
	c5 <- "c5-2"

	close(c4)
	close(c5)

	o.Wait()
	close(result)

	for i := 0; i < len(expected2); i++ {
		s := <-result
		if g, e := s, expected2[i]; g != e {
			t.Errorf("%d got %q, want %q", i, g, e)
		}
	}

	o.Close()

	n5 := runtime.NumGoroutine()
	if n1 != n5 {
		time.Sleep(10 * time.Millisecond)
		n5 := runtime.NumGoroutine()
		if n1 != n5 {
			t.Errorf("NumGoroutine %d %d", n1, n5)
		}
	}
}

func ExampleOchan() {
	result := make(chan string, 100)
	o := NewOchan(result, 100)

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

func ExampleOchanWithStruct() {
	type X struct {
		No      int
		Message string
	}
	result := make(chan X, 100)
	o := NewOchan(result, 100)

	c1 := o.GetCh()
	c2 := o.GetCh()

	c1 <- X{No: 1, Message: "Hello c1"}
	c2 <- X{No: 2, Message: "Hello c2"}
	c1 <- X{No: 3, Message: "Bye c1"}
	c2 <- X{No: 4, Message: "Bye c2"}

	close(c1)
	close(c2)

	o.Wait()
	close(result)

	for s := range result {
		fmt.Println(s.No, s.Message)
		// Output:
		// 1 Hello c1
		// 3 Bye c1
		// 2 Hello c2
		// 4 Bye c2
	}
}
