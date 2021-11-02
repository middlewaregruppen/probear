package network

import (
	"context"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DNSProbe performs DNS lookups using the net.DefaultResolver
type DNSProbe struct {
    // The host to lookup. If this value is an IPv4 address a reverse lookup will be performed instead
	Host      string `json:"host"`
    // Timeout in seconds
	Timeout int    `json:"timeout" default:"60"`
	// interval in seconds that the resource should be probed after the last probe.
	Interval int `json:"interval"` // interval in seconds that the resource should be probed after the last probe.
	// labels to add to the metric
	StdLabels
}

var (
	dnslookup_time = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_dnslookup_time",
		Help:    "Time in milliseconds it takes to retrive the document",
		Buckets: []float64{0.1, 1, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000, 100000},
	}, []string{"name", "node", "region", "zone"})

	dnslookup_error = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_dnslookup_connection_error",
		Help: "Error trying to retive url",
	}, []string{"name", "node", "region", "zone"})
)

func (p *DNSProbe) Start() {
	// Add 0 to make sure it is exposed.
	l := prometheus.Labels{"name": p.Name, "node": p.Node, "region": p.Region, "zone": p.Zone}

	if p.Interval < 1 {
		p.Interval = 10
	}

	dnslookup_error.With(l).Add(0)

	r := net.DefaultResolver
	r.PreferGo = false

	go func() {
		for {
			start := time.Now()

            // Figure out if host is a hostname or an IPv4 address
            ip := net.ParseIP(p.Host)
            if ip != nil {
                _, err := r.LookupAddr(context.Background(), p.Host)
		    	if err != nil {
			        dnslookup_error.With(l).Inc()
			    }
            } else {
                _, err := r.LookupHost(context.Background(), p.Host)
		    	if err != nil {
			    	dnslookup_error.With(l).Inc()
			    }
            }

			duration := time.Since(start)
			dnslookup_time.With(l).Observe(float64(duration.Milliseconds()))

			time.Sleep(time.Second * time.Duration(p.Interval))
		}

	}()

}
