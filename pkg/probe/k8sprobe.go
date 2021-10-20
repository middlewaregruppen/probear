package probe

import (
	"log"
	"os"
	"time"

	"github.com/middlewaregruppen/probear/pkg/k8s"
	"github.com/middlewaregruppen/probear/pkg/network"
)

type K8SProbes struct {
	Probes Probes
}

func (p *K8SProbes) Scan(interval int) {
	go func() {
		p.PopulateProbearPods()
		time.Sleep(time.Second * time.Duration(interval))
	}()
}

func (p *K8SProbes) PopulateProbearPods() {

	pods, err := k8s.GetProbearPods()

	if err != nil {
		log.Printf("Error getting Probear pods in cluster %s ", err)
		return
	}

	hn, _ := os.Hostname()

	// this pod
	thisPod, err := k8s.GetPod(hn)
	if err != nil {
		log.Printf("Error getting Probear pods in cluster %s ", err)
		return
	}

	// Set up TCP Session Probe to communicate with all probear pods except itself.
	for _, dst := range pods {
		if dst.Name == thisPod.Name {
			continue
		}

		// Check if it already exists..

		if p.Probes.HasTCPSession(thisPod.Name) {
			continue
		}

		new := network.TCPSessionProbe{
			Addr: dst.Addr,
		}
		new.Name = dst.Name
		new.Zone = thisPod.Zone
		new.Region = thisPod.Region
		new.Node = thisPod.Node
		new.DestinationNode = dst.Node
		new.DestinationRegion = dst.Region
		new.DestinationZone = dst.Zone

		p.Probes.TCPSession = append(p.Probes.TCPSession, new)
		new.Start()

	}

}
