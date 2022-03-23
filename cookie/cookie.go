package cookie

import "net/http"

var (
	cookieDomain = ""
	cookieSecure = false
)

func SetHTTPOnly(w http.ResponseWriter, k, v string) {
	setCookie(w, k, v, opts{HttpOnly: true})
}

func DeleteHTTPOnly(w http.ResponseWriter, k string) {
	setCookie(w, k, "", opts{Remove: true, HttpOnly: true})
}

func Set(w http.ResponseWriter, k, v string) {
	setCookie(w, k, v, opts{})
}

func Delete(w http.ResponseWriter, k string) {
	setCookie(w, k, "", opts{Remove: true})
}

func SetDomain(d string) {
	cookieDomain = d
}

func SetSecure(s bool) {
	cookieSecure = s
}

type opts struct {
	Remove   bool
	HttpOnly bool
}

func setCookie(w http.ResponseWriter, name, value string, opts opts) {
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

	c.Domain = cookieDomain

	if cookieSecure {
		c.Secure = true
	}

	http.SetCookie(w, &c)
}
