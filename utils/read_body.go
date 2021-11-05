package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func ReadBody(res *fasthttp.Response, logger *logrus.Logger) []byte {
	var body []byte
	var err error
	switch string(res.Header.Peek("content-encoding")) {
	case "br":
		body, err = res.BodyUnbrotli()
	case "gzip":
		body, err = res.BodyGunzip()
	case "deflate":
		body, err = res.BodyInflate()
	default:
		body = res.Body()
	}

	if err != nil {
		logger.Error(err)
	}

	return body
}
