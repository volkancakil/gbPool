package utils

import (
	"github.com/valyala/fasthttp"
)

func ReadBody(res *fasthttp.Response) ([]byte, error) {
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
		return nil, err
	}

	return body, nil
}
