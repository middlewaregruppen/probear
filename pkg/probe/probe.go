package probe

import (
	"github.com/middlewaregruppen/probear/pkg/network"
)

type Probes struct {
	HTTPGet    []network.HTTPGetProbe    `json:"httpGetProbes,omitempty"`
	TCPConnect []network.TCPConnectProbe `json:"tcpConnectProbes,omitempty"`
	TCPSession []network.TCPSessionProbe `json:"tcpSessionProbes,omitempty"`
}

func (p *Probes) HasTCPSession(name string) bool {
	for _, v := range p.TCPSession {
		if v.Name == name {
			return true
		}
	}
	return false

}

func (p *Probes) Start() {

	// RunEachProbe.
	for _, p := range p.HTTPGet {
		p.Start()
	}
	for _, p := range p.TCPConnect {
		go p.Run()
	}
	for _, p := range p.TCPSession {
		p.Start()
	}

}
