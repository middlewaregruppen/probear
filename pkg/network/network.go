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
