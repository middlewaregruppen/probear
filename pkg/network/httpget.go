package network

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/* HTTPGet retrives a document and meassures the time it takes and if there are any errors.
 */

type HTTPGetProbe struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout" default:"60"`
	// interval in seconds that the resource should be probed after the last probe.
	Interval int `json:"interval"`
	// labels to add to the metric
	StdLabels
}

var (
	httpget_time = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_httpget_time",
		Help:    "Time in milliseconds it takes to retrive the document",
		Buckets: []float64{0.1, 1, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000, 100000},
	}, []string{"name", "node", "region", "zone"})

	httpget_error = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_httpget_connection_error",
		Help: "Error trying to retive url",
	}, []string{"name", "node", "region", "zone"})
)

func (p *HTTPGetProbe) Start() {
	// Add 0 to make sure it is exposed.
	l := prometheus.Labels{"name": p.Name, "node": p.Node, "region": p.Region, "zone": p.Zone}

	if p.Interval < 1 {
		p.Interval = 10
	}

	httpget_error.With(l).Add(0)

	go func() {
		for {
			d, _, _, err := HTTPGet(p.URL, p.Timeout)
			if err != nil {
				httpget_error.With(l).Inc()
			}
			httpget_time.With(l).Observe(float64(d.Milliseconds()))

			time.Sleep(time.Second * time.Duration(p.Interval))
		}

	}()

}

/* HTTPGet does an http get to the url and reads all of the data.
   Returns the duration, http status code, response size in bytes and error.
*/
func HTTPGet(url string, timeout int) (time.Duration, int, int, error) {

	client := http.Client{
		Timeout: time.Duration(timeout * int(time.Second)),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		duration := time.Since(start)
		return duration, 0, 0, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		duration := time.Since(start)
		return duration, 0, 0, err
	}

	duration := time.Since(start)
	return duration, resp.StatusCode, len(b), err
}
