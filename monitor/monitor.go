package monitor

import (
	"gbPool/fetcher"
	"gbPool/public"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func NewMonitor(source chan *public.Proxy, fetcher fetcher.Fetcher, logger *logrus.Logger) *Monitor {
	m := &Monitor{
		fetcher: fetcher,
		logger:  logger,
		source:  source,
	}

	go m.monitor()
	return m
}

type Monitor struct {
	fetcher  fetcher.Fetcher
	logger   *logrus.Logger
	source   chan *public.Proxy
	fetching bool
}

func (m *Monitor) monitor() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		if len(m.source) < viper.GetInt("size")*3 && !m.fetching {
			m.fetching = true
			go m.update()
		}
		<-ticker.C
	}
}

func (m *Monitor) update() {
	m.fetcher.Do()
	m.fetching = false
}
