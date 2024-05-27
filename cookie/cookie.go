package cookie

import (
	"net/http"
	"regexp"
)

var (
	cookieDomainRegex = &regexp.Regexp{}
	cookieSecure      = false
	refererHeaderKey  = "Referer"
)

type Opts struct {
	Domain string
}

func SetHTTPOnly(r *http.Request, w http.ResponseWriter, k, v string, o ...Opts) {
	inner := opts{HttpOnly: true}
	if outer := getOpts(o...); outer.Domain != "" {
		inner.Domain = outer.Domain
	}
	setCookie(r, w, k, v, inner)
}

func DeleteHTTPOnly(r *http.Request, w http.ResponseWriter, k string, o ...Opts) {
	inner := opts{Remove: true, HttpOnly: true}
	if outer := getOpts(o...); outer.Domain != "" {
		inner.Domain = outer.Domain
	}
	setCookie(r, w, k, "", inner)
}

func Set(r *http.Request, w http.ResponseWriter, k, v string, o ...Opts) {
	inner := opts{}
	if outer := getOpts(o...); outer.Domain != "" {
		inner.Domain = outer.Domain
	}
	setCookie(r, w, k, v, inner)
}

func Delete(r *http.Request, w http.ResponseWriter, k string, o ...Opts) {
	inner := opts{Remove: true}
	if outer := getOpts(o...); outer.Domain != "" {
		inner.Domain = outer.Domain
	}
	setCookie(r, w, k, "", inner)
}

func SetDomain(r *regexp.Regexp) {
	cookieDomainRegex = r
}

func SetSecure(s bool) {
	cookieSecure = s
}

type opts struct {
	Remove   bool
	HttpOnly bool
	Domain   string
}

func setCookie(r *http.Request, w http.ResponseWriter, name, value string, opts opts) {
	c := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	if opts.HttpOnly {
		c.HttpOnly = true
	}

	if opts.Remove {
		c.MaxAge = -1
	} else {
		c.MaxAge = 3 * 60 * 60
	}

	if r != nil && cookieDomainRegex.MatchString(r.Header.Get(refererHeaderKey)) {
		c.Domain = cookieDomainRegex.FindString(r.Header.Get(refererHeaderKey))
	}

	if opts.Domain != "" {
		c.Domain = opts.Domain
	}

	if cookieSecure {
		c.Secure = true
	}

	http.SetCookie(w, &c)
}

func getOpts(o ...Opts) Opts {
	if len(o) == 1 {
		return o[0]
	}
	return Opts{}
}
