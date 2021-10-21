package network

import (
	"log"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type TCPConnectProbe struct {
	Addr     string `json:"addr"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
	StdLabels
}

var (
	tcpconnect_time = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_tcpconnect_time",
		Help:    "Probear tcp connect time in milliseconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000},
	}, []string{"name", "node", "region", "zone"})

	tcp_connect_failed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_tcpconnect_failed_connections",
		Help: "Number of connections failed",
	}, []string{"name", "node", "region", "zone"})
)

func (tc *TCPConnectProbe) Run() {

	if tc.Interval < 1 {
		tc.Interval = 10
	}

	l := prometheus.Labels{"name": tc.Name, "node": tc.Node, "region": tc.Region, "zone": tc.Zone}

	tcp_connect_failed.With(l).Add(0)
	for {

		d, err := TCPConnect(tc.Addr, tc.Timeout)

		if err != nil {
			tcp_connect_failed.With(l).Inc()
		}
		tcpconnect_time.With(l).Observe(float64(d.Milliseconds()))
		time.Sleep(time.Second * time.Duration(tc.Interval))
	}
}

func TCPConnect(addr string, timeout int) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout*int(time.Second)))
	if err != nil {
		log.Printf("TCPConnectError to %s: %s", addr, err)
		duration := time.Since(start)
		return duration, err
	}
	defer conn.Close()
	duration := time.Since(start)
	return duration, nil
}
