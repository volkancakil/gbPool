package gbPool

import (
	"fmt"
	"time"
)

// NewMonitor Create a new proxy monitor, Monitor is used to monitor the numbers of the proxy channel
// if its proxies is fewer than edge, re-fill it.
// source is the proxy channel that needs to be monitored, fetcher is the fetcher used when monitor needs to re-fill channel
// edge is the "danger" low limit of the channel.
func NewMonitor(source chan *Proxy, fetcher Fetcher, edge int) (*Monitor, error) {
	m := &Monitor{
		fetcher: fetcher,
		source:  source,
		edge: edge,
	}

	go m.monitor()
	return m, nil
}

type Monitor struct {
	fetcher Fetcher
	source  chan *Proxy
	fetching bool
	edge	 int
}

func (m *Monitor) monitor() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		if len(m.source) < m.edge && !m.fetching {
			m.fetching = true
			go m.update()
		}
		<-ticker.C
	}
}

func (m *Monitor) update() {
	err := m.fetcher.Do()
	if err != nil {
		// TODO: Is any other good way for exposing error to client?
		fmt.Printf("Error: %s\n", err)
		return
	}
	m.fetching = false
}
