package gbPool

// NewProxyPool Create a proxy pool. Proxy pool is used to contain all kinds of manager,
// and all manager share a single proxy channel.
// channelSize is the size of the proxy channel.Strongly suggest set a larger num than you really needs for buffering
func NewProxyPool(channelSize int) *ProxyPool {
	return &ProxyPool{
		ProxyChan: make(chan *Proxy, channelSize),
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
	ProxyChan chan *Proxy
	ProxyMgr  map[string]*manager
}
