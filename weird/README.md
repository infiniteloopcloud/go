# Weird error

Yet another error management tool

### Usage

```go
package main

import (
	"errors"
	
	"gitlab.com/infiniteloopcloud/weird"
)

var ErrTest = errors.New("test error")

func main() {
	err := weird.New("Message returned over HTTP", ErrTest, 402)
	
	errors.Is(err, ErrTest)
	// returns true
}
```