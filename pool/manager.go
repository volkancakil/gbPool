package pool

import (
	"errors"
	"fmt"
	"gbPool/fetcher"
	"gbPool/monitor"
	"gbPool/public"
	"strconv"
	"sync"
)

type manager struct {
	sync.Mutex
	enable    bool
	monitor   *monitor.Monitor
	fetcher   fetcher.Fetcher
	proxyChan chan *public.Proxy
}

// NewManager Create a new manager for a provider, Manager is used to handle fetch proxy & monitor proxy status
// managerType is the type of manager needs to be created.
// config is the specific type of config for that manager
func (p *ProxyPool) NewManager(managerType string, config interface{}) error {
	var err error
	if p.ProxyMgr[managerType] == nil {
		p.ProxyMgr[managerType], err = newProxyManager(managerType, p.ProxyChan, config)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("%s manager already exist", managerType))
	}
}

func newProxyManager(managerType string, proxyChan chan *public.Proxy, config interface{}) (*manager, error) {
	edge := 0
	mgr := &manager{}
	mgr.proxyChan = proxyChan

	switch managerType {
	case "ihuan":
		c := config.(*public.IhuanConfig)

		// If conversion failed, edge will be 0, and that's normal and ok
		edge, _ = strconv.Atoi(c.Num)
		f := fetcher.NewIhuanFetcher(mgr.proxyChan, c)
		mgr.fetcher = f
	default:
		return nil, errors.New(fmt.Sprintf("unknown type: %s", managerType))
	}

	if mgr.fetcher == nil {
		return nil, errors.New(fmt.Sprintf("%s fetcher is nil", managerType))
	}

	m, err := monitor.NewMonitor(mgr.proxyChan, mgr.fetcher, edge)
	if err != nil {
		return nil, err
	}
	mgr.monitor = m
	if mgr.monitor == nil {
		return nil, errors.New(fmt.Sprintf("%s monitor is nil", managerType))
	}

	return mgr, nil
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
