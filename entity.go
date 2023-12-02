package masscan

type MasscanResult struct {
	Hosts []Hosts `json:"hosts"`
	Ports []Ports `json:"ports"`
}

type Hosts struct {
	IP        string  `json:"ip"`
	Ports     []Ports `json:"ports"`
	Timestamp string  `json:"timestamp"`
}

type Ports struct {
	Port   int    `json:"port"`
	Proto  string `json:"proto"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	TTL    int    `json:"ttl"`
}
