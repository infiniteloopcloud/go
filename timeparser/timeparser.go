package timeparser

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

const (
	creditCardExpirationRFC3339     = "^([0-9]{4}-[0-9]{2}-[0-9]{2}[Tt][0-9]{2}:[0-9]{2}:[0-9]{2}[Zz+-:0-9]{1,6}$)"
	creditCardExpirationRFC3339Nano = "^([0-9]{4}-[0-9]{2}-[0-9]{2}[Tt][0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{7,9}[Zz+-:0-9]{1,6}$)"
	creditCardExpirationRFC1123     = "^([A-Za-z]{3}, [0-9]{2} [A-Za-z]{3} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Za-z]{3,4}$)"
	creditCardExpirationRFC1123Z    = "^([A-Za-z]{3}, [0-9]{2} [A-Za-z]{3} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} [-+]{1}[0-9]{4}$)"
	creditCardExpirationRFC822Z     = "^([0-9]{2} [A-Za-z]{3} [0-9]{2} [0-9]{2}:[0-9]{2} [-+]{1}[0-9]{4}$)"
	creditCardExpirationRFC822      = "^([0-9]{2} [A-Za-z]{3} [0-9]{2} [0-9]{2}:[0-9]{2} [A-Za-z]{3,4}$)"
	creditCardExpirationRFC850      = "^([A-Za-z]{6,9}, [0-9]{2}-[A-Za-z]{3}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Z]{3,4}$)"
	creditCardExpirationRubyFormat  = "^([A-Za-z]{3} [A-Za-z]{3} [0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} [+-][0-9]{4} [0-9]{4}$)"
	creditCardExpirationUnixFormat  = "^([A-Za-z]{3} [A-Za-z]{3} [0-9_ ]{1,2} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Za-z]{3,4} [0-9]{4}$)"
	creditCardExpirationANSICFormat = "^([A-Za-z]{3} [A-Za-z]{3} [0-9_ ]{1,2} [0-9]{2}:[0-9]{2}:[0-9]{2} [0-9]{4}$)"
)

var (
	CreditCardExpirationRFC3339Regex     = regexp.MustCompile(creditCardExpirationRFC3339)
	CreditCardExpirationRFC3339NanoRegex = regexp.MustCompile(creditCardExpirationRFC3339Nano)
	CreditCardExpirationRFC1123Regex     = regexp.MustCompile(creditCardExpirationRFC1123)
	CreditCardExpirationRFC1123ZRegex    = regexp.MustCompile(creditCardExpirationRFC1123Z)
	CreditCardExpirationRFC822ZRegex     = regexp.MustCompile(creditCardExpirationRFC822Z)
	CreditCardExpirationRFC822Regex      = regexp.MustCompile(creditCardExpirationRFC822)
	CreditCardExpirationRFC850Regex      = regexp.MustCompile(creditCardExpirationRFC850)
	CreditCardExpirationRubyFormatRegex  = regexp.MustCompile(creditCardExpirationRubyFormat)
	CreditCardExpirationUnixFormatRegex  = regexp.MustCompile(creditCardExpirationUnixFormat)
	CreditCardExpirationANSICFormatRegex = regexp.MustCompile(creditCardExpirationANSICFormat)
)

var (
	ErrUnknownFormat = errors.New("unknown format of expiration date")
	ErrInvalidMonth  = errors.New("invalid month in expiration date")
)

func Get(s string) (time.Time, error) {
	switch {
	case CreditCardExpirationRFC3339Regex.MatchString(s):
		return timeParser(time.RFC3339, s)
	case CreditCardExpirationRFC3339NanoRegex.MatchString(s):
		return timeParser(time.RFC3339Nano, s)
	case CreditCardExpirationRFC1123ZRegex.MatchString(s):
		return timeParser(time.RFC1123Z, s)
	case CreditCardExpirationRFC1123Regex.MatchString(s):
		return timeParser(time.RFC1123, s)
	case CreditCardExpirationRFC850Regex.MatchString(s):
		return timeParser(time.RFC850, s)
	case CreditCardExpirationRFC822Regex.MatchString(s):
		return timeParser(time.RFC822, s)
	case CreditCardExpirationRFC822ZRegex.MatchString(s):
		return timeParser(time.RFC822Z, s)
	case CreditCardExpirationRubyFormatRegex.MatchString(s):
		return timeParser(time.RubyDate, s)
	case CreditCardExpirationUnixFormatRegex.MatchString(s):
		return timeParser(time.UnixDate, s)
	case CreditCardExpirationANSICFormatRegex.MatchString(s):
		return timeParser(time.ANSIC, s)
	}

	return time.Time{}, ErrUnknownFormat
}

func timeParser(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil && strings.Contains(err.Error(), "month out of range") {
		return t, ErrInvalidMonth
	}
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}
