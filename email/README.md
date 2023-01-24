# Email

### Usage

```go
package main

import (
	"context"

	"gitlab.com/metricsglobal/misc-go/email"
)

func main() {
	provider, err := email.NewProvider(context.Background(), email.Opts{
		Provider: email.SMTPProvider, // Use smtp config
		SmtpConfig: &email.SmtpConfig{
			Host:     "smtp.google.com",
			Port:     587,
			Username: "username",
			Password: "<password>",
		},
	})
	if err != nil {
		panic(err)
	}

	// Set the data
	data := email.Data{
		Sender:   "john_testington@mgi.com",
		Subject:  "Important email",
		BodyHTML: `<html><body><h1>Welcome to my site</h1></body></html>`,
		Body:     "Welcome to my site",
		Charset:  email.DefaultCharset,
		ToAddresses: []string{
			"jane_testington@mgi.com",
		},
		CCAddresses: nil,
	}

	_, err = provider.Send(context.Background(), data)
	if err != nil {
		panic(err)
    }
}
```
