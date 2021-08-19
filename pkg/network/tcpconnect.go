package network

import (
	"fmt"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	tcpconnect_time = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_tcpconnect_time",
		Help:    "Probear tcp connect time in milliseconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000, 100000},
	}, []string{"probe"})

	tcp_connect_failed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_tcpconnect_failed_connections",
		Help: "Number of connections failed",
	}, []string{"probe"})
)

type TCPConnectProbe struct {
	Addr    string            `json:"addr"`
	Timeout int               `json:"timeout"`
	Status  *TCPConnectStatus `json:"status"`
	Labels  prometheus.Labels `json:"-"`
}

type TCPConnectStatus struct {
	Status
}

func (tc *TCPConnectProbe) Probe() {
	d, err := TCPConnect(tc.Addr, tc.Timeout)

	tcp_connect_failed.With(tc.Labels).Add(0)

	tc.Status = &TCPConnectStatus{}
	tc.Status.ProbedAt = time.Now()
	tc.Status.Duration = d
	if err != nil {
		tcp_connect_failed.With(tc.Labels).Inc()
		tc.Status.Error = fmt.Sprintf("%s", err)
	}
	tcpconnect_time.With(tc.Labels).Observe(float64(d.Milliseconds()))
}

func TCPConnect(addr string, timeout int) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout*int(time.Second)))
	if err != nil {
		duration := time.Since(start)
		return duration, err
	}
	defer conn.Close()
	duration := time.Since(start)
	return duration, nil
}
