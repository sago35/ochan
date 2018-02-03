# Ochan

作成した順に出力するchan

## Usage

```go
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
```

## Licence

MIT
