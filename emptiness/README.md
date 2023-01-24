# Emptiness

A library with [zero](https://github.com/golang/go/issues/5901) detection helpers.

### Usage

```go
package main

import "gitlab.com/metricsglobal/misc-go/emptiness"

var _ emptiness.IsZeroer = Dataset{}

type Dataset struct {
	Field1 string
	Field2 uint
}

func (d Dataset) IsZero() bool {
	return emptiness.IsZero(d)
}
```