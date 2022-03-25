# fireplace

Library implemented to replace a line in a file. The main goal behind the project is to work with very large files as well. Tested for 4.000.000 million lines.

### Usage

```go
package main

import (
	"log"

	"github.com/infiniteloopcloud/fireplace"
)

func main() {
	// Create new instance
	h := fireplace.New()

	// Add files to work with
	err := h.Add("file1", "./fixtures/file1")
	if err != nil {
		log.Print(err)
	}
	err = h.Add("file2", "./fixtures/file2")
	if err != nil {
		log.Print(err)
	}

	// Determine what to look for
	var lookFor = []string{
		"A75AA268-4738-9C61-A132-4A3105E9440B",
		"52378C39-4E9B-9400-058A-7B76C8F4C710",
		"6C662E3C-DFD6-94C6-2C13-72778A1535BE",
		"AB9802F1-283A-EBAA-52ED-8A096D946D1F",
	}
	// Initialize the replaces
	var param = make(map[string][]*fireplace.ReplaceOpts)

	// Look for lines and setup replace options
	for _, line := range lookFor {
		result, err := h.Search(line)
		if err != nil {
			log.Print(err)
		}
		param[result.FileID] = append(param[result.FileID], &fireplace.ReplaceOpts{
			LineNumber: result.LineNumber,
			ReplaceTo:  result.LineContent + " replaced", // TODO here we need to determine what to change
		})
	}

	// Replace multiple lines
	err = h.ReplaceMultiple(param)
	if err != nil {
		log.Print(err)
	}
}
```
