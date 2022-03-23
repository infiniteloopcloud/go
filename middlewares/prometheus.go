package middlewares

import (
	"log"
	"net/http"
	"strconv"

	"github.com/infiniteloopcloud/hyper"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// RequestTotalCount counter collector stores the request total count
	requestTotalCount = "requestTotalCount"

	// RequestCounter counter collector stores the request count per request path
	requestCounter = "requestCounter"

	// HttpDuration histogram collector stores the request durations per request path
	httpDuration = "httpDuration"

	// ResponseStatus histogram collector stores the request statuses per request path
	responseStatus = "responseStatus"

	// labels
	pathLabel       = "path"
	methodLabel     = "method"
	statusCodeLabel = "status_code"
)

func init() {
	createCollectors()
	registerCollectors()
}

func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI

		timer := newTimer(histogramVec(httpDuration).WithLabelValues(path, r.Method))
		defer timer.ObserveDuration()

		next.ServeHTTP(w, r)

		counterVec(responseStatus).WithLabelValues(path, r.Method, strconv.Itoa(w.(*hyper.Writer).StatusCode)).Inc()
		counterVec(requestTotalCount).WithLabelValues("/api/v1").Inc()
		counterVec(requestCounter).WithLabelValues(path, r.Method).Inc()
	})
}

type collectors struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

// c stores all the collectors we want to use in core
var c collectors

// createCollectors creates the collectors that we are using
func createCollectors() {
	c = collectors{
		counters: map[string]*prometheus.CounterVec{
			requestTotalCount: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Number of requests",
				}, []string{pathLabel}),
			requestCounter: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests",
					Help: "Number of requests per path",
				}, []string{pathLabel, methodLabel}),
			responseStatus: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "response_status",
					Help: "Status of HTTP response",
				},
				[]string{pathLabel, methodLabel, statusCodeLabel},
			),
		},
		histograms: map[string]*prometheus.HistogramVec{
			httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: "http_response_time_seconds",
				Help: "Duration of HTTP requests.",
				// TODO define proper buckets for response_time_seconds
			}, []string{pathLabel, methodLabel}),
		},
	}
}

// registerCollectors registers the provided prometheus.Collectors with the DefaultRegisterer.
// AlreadyRegisteredError is returned by the Register method if the Collector to
// be registered has already been registered before, or a different Collector
// that collects the same metrics has been registered before.
func registerCollectors() {
	for cName, counterVec := range c.counters {
		if err := prometheus.Register(counterVec); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				if c.counters[cName], ok = are.ExistingCollector.(*prometheus.CounterVec); !ok {
					log.Println("unable to cast ExistingCollector to *prometheus.CounterVec")
				}
			} else {
				log.Println(err)
			}
		}
	}
	for hName, histogramVec := range c.histograms {
		if err := prometheus.Register(histogramVec); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				if c.histograms[hName], ok = are.ExistingCollector.(*prometheus.HistogramVec); !ok {
					log.Println("unable to cast ExistingCollector to *prometheus.HistogramVec")
				}
			} else {
				log.Println(err)
			}
		}
	}
}

// CounterVec returns the corresponding registered prometheus.CounterVec
func counterVec(name string) *prometheus.CounterVec {
	return c.counters[name]
}

// HistogramVec returns the corresponding registered prometheus.HistogramVec
func histogramVec(name string) *prometheus.HistogramVec {
	return c.histograms[name]
}

// NewTimer creates a new prometheus.Timer. The provided Observer is used to observe a
// duration in seconds. Timer is usually used to time a function call in the
// following way:
// t := prometheus.NewTimer
// defer t.ObserveDuration
func newTimer(o prometheus.Observer) *prometheus.Timer {
	return prometheus.NewTimer(o)
}
