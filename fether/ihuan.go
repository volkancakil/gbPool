package fether

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
	"github.com/jobber2955/gbPool/public"
	"github.com/jobber2955/gbPool/utils"
	"github.com/valyala/fasthttp"
	"regexp"
	"strings"
	"time"
)

var (
	ihuanPage  = "https://ip.ihuan.me/ti.html"
	ihuanApi   = "https://ip.ihuan.me/tqdl.html"
	ihuanCache = "https://ip.ihuan.me/mouse.do"
)

// NewIhuanFetcher Create a new ihuan proxy fetcher, dest is the destination proxy channel, where fetched proxies are going.
// config is the specific config struct for ihuan
func NewIhuanFetcher(dest chan *public.Proxy, config *public.IHuanConfig) *ihuanFetcher {
	client := &fasthttp.Client{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &ihuanFetcher{
		httpClient: client,
		dest:       dest,
		config: config,
	}
}

type ihuanFetcher struct {
	httpClient *fasthttp.Client
	dest       chan *public.Proxy
	key        string
	rawIps     string
	config	   *public.IHuanConfig
}

func (i *ihuanFetcher) Do() error {
	if i.key == "" {
		if err := i.init(); err != nil {
			return err
		}
	}
	if err := i.fetch(); err != nil {
		return err
	}
	i.parse()
	return nil
}

func (i *ihuanFetcher) init() error {
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.SetRequestURI(ihuanPage)
	req.Header.HasAcceptEncoding("gzip, deflate, br")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	if err := i.httpClient.Do(req, res); err != nil {
		return err
	}

	if res.Header.Peek("set-cookie") != nil {
		cookie := string(res.Header.Peek("set-cookie"))
		var err error

		if err = i.loadCache(cookie); err != nil {
			return err
		}
	} else {
		body, err := utils.ReadBody(res)
		if err != nil {
			return err
		}
		if doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
			return err
		} else {
			i.key = doc.Find(`[name="key"]`).AttrOr("value", "")
		}
	}
	return nil
}

func (i *ihuanFetcher) loadCache(cookie string) error {
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.SetRequestURI(ihuanCache)
	req.Header.HasAcceptEncoding("gzip, deflate, br")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", "https://ip.ihuan.me/ti.html")
	if err := i.httpClient.Do(req, res); err != nil {
		return err
	}

	body, err := utils.ReadBody(res)
	if err != nil {
		return err
	}

	rePattern, err := regexp.Compile("[a-zA-z0-9]{32}")
	if err != nil {
		return err
	}

	cache := rePattern.Find(body)
	if cache != nil {
		i.key = string(cache)
		return  nil
	} else {
		return errors.New(fmt.Sprintf("failed to find cache key, raw: %s", string(body)))
	}
}

func (i *ihuanFetcher) fetch() error {
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	postStruct := &ihuanPost{
		Num:         i.config.Num,
		Anonymity:   i.config.Anonymity,
		Type:        i.config.Type,
		Post:        i.config.Post,
		Sort:        i.config.Sort,
		Port:        i.config.Port,
		KillPort:    i.config.KillPort,
		Address:     i.config.Address,
		KillAddress: i.config.KillAddress,
		Key:         i.key,
	}
	postValues, err := query.Values(postStruct)
	if err != nil {
		return err
	}

	req.SetRequestURI(ihuanApi)
	req.SetBodyString(postValues.Encode())
	req.Header.SetReferer("https://ip.ihuan.me/ti.html")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	if err := i.httpClient.Do(req, res); err != nil {
		return err
	}

	body, err := utils.ReadBody(res)
	if err != nil {
		return err
	}

	if doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
		return err
	} else {
		rawIPs, err := doc.Find(`div.col-md-10>div.panel.panel-default>div.panel-body`).Html()
		if err != nil {
			return err
		}

		i.rawIps = rawIPs
		return nil
	}
}

func (i *ihuanFetcher) parse() {
	for _, ip := range strings.Split(i.rawIps, "<br/>") {
		i.dest <- &public.Proxy{
			Address: ip,
			Expire:  0,
		}
	}
}

type ihuanPost struct {
	Num         string `url:"num"`          // ????????????
	Anonymity   string `url:"anonymity"`    // ???????????? 0)???????????? 1)???????????? 2)???????????? NIL)????????????
	Type        string `url:"type"`         // ???????????? 0)???HTTP 1)???HTTPS NIL)??????
	Post        string `url:"post"`         // ???????????? 1)??????POST NIL)??????
	Sort        string `url:"sort"`         // ???????????? 1)???????????????????????? 2)???????????????????????? 3)???????????????????????? 4)???????????????????????? NIL)??????
	Port        string `url:"port"`         // ????????????
	KillPort    string `url:"kill_port"`    // ????????????
	Address     string `url:"address"`      // ????????????
	Key         string `url:"key"`          // ??????
	KillAddress string `url:"kill_address"` // ????????????
}
