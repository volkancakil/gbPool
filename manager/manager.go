package manager

import (
	"errors"
	"fmt"
	"github.com/jobber2955/gbPool/fether"
	monitor2 "github.com/jobber2955/gbPool/monitor"
	"github.com/jobber2955/gbPool/public"
	"strconv"
	"sync"
)

type Manager struct {
	sync.Mutex
	enable    bool
	monitor   *monitor2.Monitor
	fetcher   fether.Fetcher
	proxyChan chan *public.Proxy
}

func NewProxyManager(managerType string, proxyChan chan *public.Proxy, config interface{}) (*Manager, error) {
	edge := 0
	mgr := &Manager{}
	mgr.proxyChan = proxyChan

	switch managerType {
	case "ihuan":
		c := config.(*public.IHuanConfig)

		// If conversion failed, edge will be 0, and that's normal and ok
		edge, _ = strconv.Atoi(c.Num)
		f := fether.NewIhuanFetcher(mgr.proxyChan, c)
		mgr.fetcher = f
	default:
		return nil, errors.New(fmt.Sprintf("unknown type: %s", managerType))
	}

	if mgr.fetcher == nil {
		return nil, errors.New(fmt.Sprintf("%s fetcher is nil", managerType))
	}

	m, err := monitor2.NewMonitor(mgr.proxyChan, mgr.fetcher, edge)
	if err != nil {
		return nil, err
	}
	mgr.monitor = m
	if mgr.monitor == nil {
		return nil, errors.New(fmt.Sprintf("%s monitor is nil", managerType))
	}

	return mgr, nil
}

func (m *Manager) Enable() {
	defer m.Unlock()
	m.Lock()
	m.enable = true
}

func (m *Manager) Disable() {
	defer m.Unlock()
	m.Lock()
	for len(m.proxyChan) != 0 {
		<-m.proxyChan
	}
	m.enable = false
}
