package pool

import (
	"gbPool/public"
)

func NewProxyPool(channelSize int) *ProxyPool {
	return &ProxyPool{
		ProxyChan: make(chan *public.Proxy, channelSize),
		ProxyMgr:  map[string]*manager{},
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
	ProxyMgr  map[string]*manager
}
