package network

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	tcpsession_time = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "probear_tcpsession_pingpong_time",
		Help:    "Time it takes to send a message to the server and have the message sent back from the server in milliseconds (ms)",
		Buckets: []float64{0.0000001, 0.000001, 0.00001, 0.0001, 0.001, 0.01, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1000, 10000, 100000},
	}, []string{"name", "node", "region", "zone", "dstnode", "dstregion", "dstzone"})

	tcpsession_failed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_tcpsession_failed_sessions",
		Help: "Number of established sessions that has failed",
	}, []string{"name", "node", "region", "zone", "dstnode", "dstregion", "dstzone"})

	tcpsession_failed_conn = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "probear_tcpsession_failed_connections",
		Help: "Number of connections attempts that has failed",
	}, []string{"name", "node", "region", "zone", "dstnode", "dstregion", "dstzone"})
)

type TCPSessionProbe struct {
	Addr     string `json:"addr"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
	StdLabels
	DestinationZone   string `json:"destinationZone"`
	DestinationRegion string `json:"destinationRegion"`
	DestinationNode   string `json:"destinationNode"`

	started bool
}

type res struct {
	failedSession    bool
	failedConnection bool
	duration         time.Duration
	err              error
}

func (t *TCPSessionProbe) Start() {

	if !t.started {
		go t.runClient()
		t.started = true
	}

}

func (t *TCPSessionProbe) runClient() error {
	l := prometheus.Labels{
		"name":      t.Name,
		"node":      t.Node,
		"region":    t.Region,
		"zone":      t.Zone,
		"dstnode":   t.DestinationNode,
		"dstregion": t.DestinationRegion,
		"dstzone":   t.DestinationZone,
	}

	// Set vectors to 0 so they register with prometheus.
	tcpsession_failed.With(l).Add(0)
	tcpsession_failed_conn.With(l).Add(0)

	reschan := make(chan res, 10)

	if t.Interval < 1 {
		t.Interval = 1
	}
	if t.Timeout < 1 {
		t.Interval = 10
	}

	go TCPSessionClient(t.Addr, time.Second*time.Duration(t.Interval), time.Second*time.Duration(t.Timeout), reschan)

	for {

		d := <-reschan

		// Update status
		if d.failedConnection {
			tcpsession_failed_conn.With(l).Inc()
		}
		if d.failedSession {
			tcpsession_failed.With(l).Inc()
		}

		// Prometheus histogram
		if d.duration > 0 {
			tcpsession_time.With(l).Observe(float64(d.duration.Milliseconds()))
		}
	}

}

func TCPSessionClient(addr string, interval time.Duration, timeout time.Duration, ch chan res) error {

	for {
		var conn net.Conn
		var err error

		for {
			conn, err = net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				log.Printf("tcp session client: error connecting to session server %s ", err)
				ch <- res{err: err, failedConnection: true}
				time.Sleep(5 * time.Second)
				continue
			}
			defer conn.Close()
			break
		}

		for {

			start := time.Now()

			// Send ping to server
			msg := []byte("probearhello")
			_, err := conn.Write(msg)

			if err != nil {
				ch <- res{err: err, failedSession: true}
				break
			}

			// Receive ping from server.
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				ch <- res{err: err, failedSession: true}
				break
			}

			duration := time.Since(start)
			ch <- res{duration: duration}
			time.Sleep(interval)

		}
	}
}

/* 	TCPSessionServer accepts tcp connections and waits on a message.
Once a message is received it is echoed back to the client. */
func TCPSessionServer(port int) {
	laddr := fmt.Sprintf("0.0.0.0:%d", port)
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// accept connection
		log.Println("tcp session server: waiting on new connection...")
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("tcp session server: connection from %s", conn.RemoteAddr())

		// handle connection
		go handleSessionConn(conn)
	}
}

func handleSessionConn(c net.Conn) {
	defer c.Close()

	for {
		// handle incoming data
		buffer := make([]byte, 1024)
		numBytes, err := c.Read(buffer)
		if err != nil {

			log.Printf("tcp session server: connection to %s lost %s", c.RemoteAddr(), err)
			break
		}
		// handle reply
		msg := string(buffer[:numBytes])
		_, err = c.Write([]byte(msg))
		if err != nil {
			log.Printf("tcp session server: connection to %s lost %s", c.RemoteAddr(), err)
			break
		}

	}
}
