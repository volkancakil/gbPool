package pool

import (
	"gbPool/public"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.AutomaticEnv()
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetConfigName("proxy_pool")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic("Failed to read in viper config")
		os.Exit(-1)
	}
}

func NewProxyPool(logger *logrus.Logger) *ProxyPool {
	return &ProxyPool{
		ProxyChan: make(chan *public.Proxy, viper.GetInt("size")*30),
		ProxyMgr:  map[string]*manager{},
		logger:    logger,
	}
}

func (p *ProxyPool) NewManager(managerType string) {
	if p.ProxyMgr[managerType] == nil {
		p.ProxyMgr[managerType] = newProxyManager(managerType, p.logger, p.ProxyChan)
	} else {
		p.logger.Warnf("Already created: %s\n", managerType)
	}
}

func (p *ProxyPool) Switch(old, new string) {
	if p.ProxyMgr[old] != nil {
		p.ProxyMgr[old].Disable()
	}

	p.NewManager(new)
	p.ProxyMgr[new].Enable()
}

type ProxyPool struct {
	ProxyChan chan *public.Proxy
	ProxyMgr  map[string]*manager
	logger    *logrus.Logger
}
