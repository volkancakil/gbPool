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
func NewIhuanFetcher(dest chan *public.Proxy, config *public.IhuanConfig) *ihuanFetcher {
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
	config	   *public.IhuanConfig
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
	Num         string `url:"num"`          // 提取数量
	Anonymity   string `url:"anonymity"`    // 匿名程度 0)透明代理 1)普匿代理 2)高匿代理 NIL)全部代理
	Type        string `url:"type"`         // 代理类型 0)仅HTTP 1)仅HTTPS NIL)不限
	Post        string `url:"post"`         // 代理模式 1)支持POST NIL)不限
	Sort        string `url:"sort"`         // 排序方式 1)验证时间从近到远 2)验证时间从远到近 3)存活时间从短到长 4)存活时间从长到短 NIL)随机
	Port        string `url:"port"`         // 指定端口
	KillPort    string `url:"kill_port"`    // 排除端口
	Address     string `url:"address"`      // 指定地区
	Key         string `url:"key"`          // 密钥
	KillAddress string `url:"kill_address"` // 排除地区
}
