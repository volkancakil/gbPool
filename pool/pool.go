package pool

import (
	"errors"
	"fmt"
	"github.com/jobber2955/gbPool/manager"
	"github.com/jobber2955/gbPool/public"
)

// NewProxyPool Create a proxy pool. Proxy pool is used to contain all kinds of manager,
// and all manager share a single proxy channel.
// channelSize is the size of the proxy channel.Strongly suggest set a larger num than you really needs for buffering
func NewProxyPool(channelSize int) *ProxyPool {
	return &ProxyPool{
		ProxyChan: make(chan *public.Proxy, channelSize),
		ProxyMgr:  map[string]*manager.Manager{},
	}
}

// NewManager Create a new manager for a provider, Manager is used to handle fetch proxy & monitor proxy status
// managerType is the type of manager needs to be created.
// config is the specific type of config for that manager
func (p *ProxyPool) NewManager(managerType string, config interface{}) error {
	var err error
	if p.ProxyMgr[managerType] == nil {
		p.ProxyMgr[managerType], err = manager.NewProxyManager(managerType, p.ProxyChan, config)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("%s manager already exist", managerType))
	}
}

func (p *ProxyPool) Switch(old, new string) {
	if p.ProxyMgr[old] != nil {
		p.ProxyMgr[old].Disable()
	}
	p.ProxyMgr[new].Enable()
}

type ProxyPool struct {
	ProxyChan chan *public.Proxy
	ProxyMgr  map[string]*manager.Manager
}
