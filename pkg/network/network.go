package network

func init() {
	// Run the TCP Session Server.
	go TCPSessionServer(10000)
}

type StdLabels struct {
	Name      string `json:"name"`
	Node      string `json:"node"`
	Region    string `json:"region"`
	Zone      string `json:"zone"`
	isManaged bool   `json:"-"`
}

func (l *StdLabels) GetName() string {
	return l.Name
}

// internal k8s probing in the namespace.

/*
func (n *Network) updateProbearK8STargets() {
	pods, err := k8s.GetProbearPods()
	if err != nil {
		log.Printf("err getting probear pods: %s ", err)
	}

	// Add new pods.
	for _, p1 := range pods {

		if registered(n.ProbearTargets, p1.Name) {
			continue
		}

		/* hostname, _ := os.Hostname()
		thisPod, err := k8s.GetPod(hostname)

		if err != nil {
			log.Printf("Can not find the pod that we are running on. %s", err)
			continue
		}
		*(/)

		labels := prometheus.Labels{
			"probename": p1.Node,
			//"target": p1.Node,
			//"sourcenode": thisPod.Node,
		}

		n.ProbearTargets = append(n.ProbearTargets,
			NetworkTarget{
				Name: fmt.Sprintf("%s", p1.Name),
				TCPConnect: &TCPConnectProbe{
					Addr:    fmt.Sprintf("%s:2112", p1.Addr),
					Labels:  labels,
					Timeout: 10,
				},
				TCPSession: &TCPSessionProbe{
					Addr:    fmt.Sprintf("%s:10000", p1.Addr),
					Timeout: 10,
				},
			})
	}
	// Remove old pods.
	for k, t := range n.ProbearTargets {
		if !contains(pods, t.Name) {
			n.ProbearTargets = append(n.ProbearTargets[:k], n.ProbearTargets[k+1:]...)
		}
	}

}

func registered(nts []NetworkTarget, name string) bool {
	for _, t := range nts {
		if t.Name == name {
			return true
		}
	}
	return false
}

func contains(s []k8s.ProbearPod, e string) bool {
	for _, a := range s {
		if a.Name == e {
			return true
		}
	}
	return false
}
*/
