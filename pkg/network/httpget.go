package network

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpget_time = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probear_httpget_time",
		Help: "Time in milliseconds it took to retrive the content from the url",
	}, []string{"probe"})
	httpget_error = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_httpget_error",
		Help: "Error trying to retive url",
	}, []string{"probe"})
)

type HTTPGetProbe struct {
	URL     string      `json:"url"`
	Timeout int         `json:"timeout" default:"60"`
	Status  *HTTPStatus `json:"status"`
}

type HTTPStatus struct {
	Status
	ResponseCode    int `json:"responseCode"`
	BytesInResponse int `json:"bytesInResponse"`
}

func (p *HTTPGetProbe) Probe(name string) {

	d, c, bz, err := HTTPGet(p.URL, p.Timeout)

	p.Status = &HTTPStatus{
		ResponseCode:    c,
		BytesInResponse: bz,
	}
	if err != nil {
		p.Status.Error = fmt.Sprintf("%s", err)
		httpget_error.With(prometheus.Labels{"probe": name}).Inc()
	}

	p.Status.ProbedAt = time.Now()
	p.Status.Duration = d

	httpget_time.With(prometheus.Labels{"probe": name}).Set(float64(d.Microseconds()))

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
