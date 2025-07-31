package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni/v3"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_shortener",
			Help: "HTTP requests in shortener app with response code and endpoint labels",
		},
		[]string{"code", "endpoint"},
	)
	RequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_shortener",
			Help:       "duration of the http request in shortener app",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"code", "endpoint"},
	)
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		then := time.Now()

		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)

		labels := prometheus.Labels{
			"code":     strconv.Itoa(lrw.Status()),
			"endpoint": r.RequestURI,
		}

		RequestsTotal.With(labels).Inc()
		RequestDuration.With(labels).Observe(time.Since(then).Seconds())
	})
}
