# Cookie

Cookie helper functions.

```go
package main

import (
	"net/http"
	
	"github.com/infiniteloopcloud/cookie"
)

func main() {
	// Set the cookie's domain options
	cookie.SetDomain(".test.com")
	// Set the cookie's secure option
	cookie.SetSecure(true)

	// Let's say this is the response
	var w http.ResponseWriter

	// Set a NOT HTTP ONLY value
 	cookie.Set(w, "jwt", "test")
	// Delete that cookie
	cookie.Delete(w, "jwt")

	// Set a HTTP ONLY value
	cookie.SetHTTPOnly(w, "jwt", "test_http_only")
	// Delete that cookie
	cookie.DeleteHTTPOnly(w, "jwt")
}
```
