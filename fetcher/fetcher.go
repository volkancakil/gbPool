package fetcher

import (
	"gbPool/public"
	"github.com/sirupsen/logrus"
)

type Fetcher interface {
	Do()
	init()
	fetch()
	parse()
}

func NewFetcher(fetcherType string, logger *logrus.Logger, dest chan *public.Proxy) Fetcher {
	var f Fetcher
	switch fetcherType {
	case "ihuan":
		f = newIhuanFetcher(logger, dest)
	default:
		logger.Error("Unknown type: %s", fetcherType)
		return nil
	}
	return f
}
