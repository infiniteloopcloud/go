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

func SetHTTPOnly(r *http.Request, w http.ResponseWriter, k, v string) {
	setCookie(r, w, k, v, opts{HttpOnly: true})
}

func DeleteHTTPOnly(r *http.Request, w http.ResponseWriter, k string) {
	setCookie(r, w, k, "", opts{Remove: true, HttpOnly: true})
}

func Set(r *http.Request, w http.ResponseWriter, k, v string) {
	setCookie(r, w, k, v, opts{})
}

func Delete(r *http.Request, w http.ResponseWriter, k string) {
	setCookie(r, w, k, "", opts{Remove: true})
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

	if cookieSecure {
		c.Secure = true
	}

	http.SetCookie(w, &c)
}
