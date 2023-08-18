[![Actions Status](https://github.com/sago35/ochan/actions/workflows/build.yml/badge.svg)](https://github.com/sago35/Getopt-Kingpin/actions) [![GoDoc Reference](https://godoc.org/github.com/sago35/ochan?status.svg)](https://godoc.org/github.com/sago35/ochan) [![Go Report Card](https://goreportcard.com/badge/github.com/sago35/ochan)](https://goreportcard.com/report/github.com/sago35/ochan)

# Ochan

Package ochan provides ordered chan.

## Usage

Package ochan supports generics.  
Thus, it can work with any type by supplying chan T to `NewOchan()`.  
See ochan_test.go for an example of its use.  

```go
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
```

## Licence

MIT
