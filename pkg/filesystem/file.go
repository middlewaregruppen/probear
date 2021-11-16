package filesystem

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type StdLabels struct {
	Name   string `json:"name"`
	Node   string `json:"node"`
	Region string `json:"region"`
	Zone   string `json:"zone"`
}

type FilesystemProbe struct {
	StdLabels
	Path     string `json:"path"`
	FileSize int    `json:"fileSize"`
	Interval int    `json:"interval"`
}

var (
	fsp_write = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_file_write_speed",
		Help:    "File write speed in kilobytes per seconds",
		Buckets: []float64{0.1, 1, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000, 100000},
	}, []string{"name", "node", "region", "zone"})

	fsp_read = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_file_read_speed",
		Help: "File read in kilobytes per seconds",
	}, []string{"name", "node", "region", "zone"})

	fsp_error = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_file_write_errors",
		Help: "File write errors",
	}, []string{"name", "node", "region", "zone"})
)

func (p *FilesystemProbe) Start() {
	l := prometheus.Labels{"name": p.Name, "node": p.Node, "region": p.Region, "zone": p.Zone}
	fsp_error.With(l).Add(0)

	if p.Interval < 1 {
		p.Interval = 20
	}

	go func() {
		for {
			d, err := writeFile(p.Path, p.FileSize)

			if err != nil {
				fsp_error.With(l).Inc()
				log.Printf("%s", err)
			}
			fsp_write.With(l).Observe(float64(d.Milliseconds()))

			time.Sleep(time.Second * time.Duration(p.Interval))
		}

	}()

}

/*
 */

func writeFile(path string, fileSize int) (time.Duration, error) {

}
