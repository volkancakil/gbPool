package pool

import (
	"gbPool/fetcher"
	"gbPool/monitor"
	"gbPool/public"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type manager struct {
	sync.Mutex
	enable    bool
	monitor   *monitor.Monitor
	fetcher   fetcher.Fetcher
	proxyChan chan *public.Proxy
}

func newProxyManager(managerType string, logger *logrus.Logger, proxyChan chan *public.Proxy) *manager {
	mgr := &manager{}
	size := viper.GetInt("size")
	if size == 0 || size > 100 {
		size = 100
	}
	mgr.proxyChan = proxyChan

	switch managerType {
	case "ihuan":
		f := fetcher.NewFetcher(managerType, logger, mgr.proxyChan)
		mgr.fetcher = f
	default:
		logger.Error("Unknown type: %s", managerType)
		return nil
	}

	if mgr.fetcher == nil {
		return nil
	}

	m := monitor.NewMonitor(mgr.proxyChan, mgr.fetcher, logger)
	mgr.monitor = m
	if mgr.monitor == nil {
		return nil
	}

	mgr.Enable()
	return mgr
}

func (m *manager) Enable() {
	defer m.Unlock()
	m.Lock()
	m.enable = true
}

func (m *manager) Disable() {
	defer m.Unlock()
	m.Lock()
	for len(m.proxyChan) != 0 {
		<-m.proxyChan
	}
	m.enable = false
}
