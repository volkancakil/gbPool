package public

import "time"

type Proxy struct {
	Address string
	Expire  int64
}

// Drop If a proxy is unusable, drop it
func (p *Proxy) Drop() {
	return
}

// ReUse Simply re-insert the proxy to channel
func (p *Proxy) ReUse(dest chan *Proxy) {
	dest <- p
}

// Expired If a proxy expired, drop it. Some kinds of proxy don't have an expiration time, so default is 0,
// should skip 0
func (p *Proxy) Expired(dest chan *Proxy) {
	if p.Expire <= time.Now().Unix() && p.Expire != 0 {
		return
	}
	p.ReUse(dest)
}

type IHuanConfig struct {
	Num         string
	Anonymity   string
	Type        string
	Post        string
	Sort        string
	Port        string
	KillPort    string
	Address     string
	Key         string
	KillAddress string
}