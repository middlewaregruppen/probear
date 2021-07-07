package network

import (
	"fmt"
	"log"
	"time"

	"github.com/middlewaregruppen/probear/pkg/k8s"
)

type Network struct {
	NetworkTargets []NetworkTarget `json:"probes"`
	ProbearTargets []NetworkTarget `json:"probears"`
}

type NetworkTarget struct {
	Name       string           `json:"name"`
	HTTPGet    *HTTPGetProbe    `json:"httpGet,omitempty"`
	TCPConnect *TCPConnectProbe `json:"tcpConnect,omitempty"`
	TCPSession *TCPSessionProbe `json:"tcpSession,omitempty"`
}

type Status struct {
	ProbedAt time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error"`
}

func (n *Network) Probe() {
	for _, t := range n.NetworkTargets {
		t.Probe()
	}

	pods, err := k8s.GetProbearPods()
	if err != nil {
		log.Printf("err getting probear pods: %s ", err)
	}

	for _, p1 := range pods {

		n.ProbearTargets = append(n.ProbearTargets,
			NetworkTarget{
				Name: fmt.Sprintf("to %s", p1),
				TCPConnect: &TCPConnectProbe{
					Addr:    fmt.Sprintf("%s:2112", p1),
					Timeout: 10,
				},
				TCPSession: &TCPSessionProbe{
					Addr:    fmt.Sprintf("%s:10000", p1),
					Timeout: 10,
				},
			})

	}

}

func (nt *NetworkTarget) Probe() {
	if nt.HTTPGet != nil {
		nt.HTTPGet.Probe(nt.Name)
	}
	if nt.TCPConnect != nil {
		nt.TCPConnect.Probe(nt.Name)
	}
	if nt.TCPSession != nil {
		nt.TCPSession.Probe(nt.Name)
	}

}
