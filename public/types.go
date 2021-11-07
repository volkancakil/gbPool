package public

import "time"

type Proxy struct {
	Address string
	Expire  int64
}

func (p *Proxy) Drop() {
	return
}

func (p *Proxy) ReUse(dest chan *Proxy) {
	dest <- p
}

func (p *Proxy) Expired(dest chan *Proxy) {
	if p.Expire <= time.Now().Unix() && p.Expire != 0 {
		return
	}
	p.ReUse(dest)
}

type IhuanConfig struct {
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