package fetcher

import (
	"bytes"
	"gbPool/public"
	"gbPool/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"regexp"
	"strings"
	"time"
)

var (
	ihuanPage  = "https://ip.ihuan.me/ti.html"
	ihuanApi   = "https://ip.ihuan.me/tqdl.html"
	ihuanCache = "https://ip.ihuan.me/mouse.do"
)

func newIhuanFetcher(logger *logrus.Logger, dest chan *public.Proxy) *ihuanFetcher {
	client := &fasthttp.Client{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	if viper.GetBool("debug") && viper.GetString("proxy") != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(viper.GetString("proxy"))
	}

	return &ihuanFetcher{
		httpClient: client,
		logger:     logger,
		dest:       dest,
	}
}

type ihuanFetcher struct {
	httpClient *fasthttp.Client
	logger     *logrus.Logger
	dest       chan *public.Proxy
	key        string
	rawIps     string
}

func (i *ihuanFetcher) Do() {
	if i.key == "" {
		i.init()
	}
	i.fetch()
	i.parse()
}

func (i *ihuanFetcher) init() {
	i.logger.Info("Loading ihuan post key")
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.SetRequestURI(ihuanPage)
	req.Header.HasAcceptEncoding("gzip, deflate, br")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	if err := i.httpClient.Do(req, res); err != nil {
		i.logger.Error(err)
		return
	}

	if res.Header.Peek("set-cookie") != nil {
		cookie := string(res.Header.Peek("set-cookie"))
		i.key = i.loadCache(cookie)
	} else {
		body := utils.ReadBody(res, i.logger)
		if doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
			i.logger.Error(err)
			return
		} else {
			i.key = doc.Find(`[name="key"]`).AttrOr("value", "")
		}
	}
}

func (i *ihuanFetcher) loadCache(cookie string) string {
	i.logger.Info("Loading ihuan post key cache")
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	req.SetRequestURI(ihuanCache)
	req.Header.HasAcceptEncoding("gzip, deflate, br")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", "https://ip.ihuan.me/ti.html")
	if err := i.httpClient.Do(req, res); err != nil {
		i.logger.Error(err)
		return ""
	}

	body := utils.ReadBody(res, i.logger)
	rePattern, err := regexp.Compile("[a-zA-z0-9]{32}")
	if err != nil {
		i.logger.Error(err)
		return ""
	}
	cache := rePattern.Find(body)
	if cache != nil {
		i.logger.Infof("Cache Found: %s\n", string(cache))
		return string(cache)
	}
	return ""
}

func (i *ihuanFetcher) fetch() {
	i.logger.Info("Fetching ihuan proxies")
	if i.key == "" {
		i.logger.Error("ihuan key empty")
		return
	}
	req := &fasthttp.Request{}
	res := &fasthttp.Response{}

	postStruct := &ihuanPost{
		Num:         viper.GetString("size"),
		Anonymity:   viper.GetString("ihuan_anonymity"),
		Type:        viper.GetString("ihuanP_type"),
		Post:        viper.GetString("ihuan_post"),
		Sort:        viper.GetString("ihuan_sort"),
		Port:        viper.GetString("ihuan_port"),
		KillPort:    viper.GetString("ihuan_kill_port"),
		Address:     viper.GetString("ihuan_address"),
		KillAddress: viper.GetString("ihuan_kill_address"),
		Key:         i.key,
	}
	postValues, err := query.Values(postStruct)
	if err != nil {
		i.logger.Error(err)
		return
	}

	req.SetRequestURI(ihuanApi)
	req.SetBodyString(postValues.Encode())
	req.Header.SetReferer("https://ip.ihuan.me/ti.html")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	if err := i.httpClient.Do(req, res); err != nil {
		i.logger.Error(err)
		return
	}

	body := utils.ReadBody(res, i.logger)
	if doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)); err != nil {
		return
	} else {
		rawIPs, err := doc.Find(`div.col-md-10>div.panel.panel-default>div.panel-body`).Html()
		if err != nil {
			return
		}
		i.rawIps = rawIPs
	}
}

func (i *ihuanFetcher) parse() {
	i.logger.Info("Parsing ihuan raw ips")
	if i.rawIps == "" {
		i.logger.Error("ihuan raw ip empty")
		return
	}
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
