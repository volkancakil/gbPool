package monitor

import (
	"fmt"
	"gbPool/fetcher"
	"gbPool/public"
	"time"
)

func NewMonitor(source chan *public.Proxy, fetcher fetcher.Fetcher, edge int) (*Monitor, error) {
	m := &Monitor{
		fetcher: fetcher,
		source:  source,
		edge: edge,
	}

	go m.monitor()
	return m, nil
}

type Monitor struct {
	fetcher  fetcher.Fetcher
	source   chan *public.Proxy
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
