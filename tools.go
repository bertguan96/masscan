package masscan

import "encoding/json"

func ParseScanResult(content []byte) (*MasscanResult, error) {
	var m []Hosts

	err := json.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}

	var result MasscanResult
	for i := range m {
		result.Ports = append(result.Ports,
			Ports{Port: m[i].Ports[0].Port,
				Proto:  m[i].Ports[0].Proto,
				Status: m[i].Ports[0].Status,
				Reason: m[i].Ports[0].Reason,
				TTL:    m[i].Ports[0].TTL,
			})
		result.Hosts = append(result.Hosts,
			Hosts{IP: m[i].IP,
				Timestamp: m[i].Timestamp,
			})

	}

	return &result, nil
}
