# GBPool-- a simple but useful golang free proxy pool
## Intro([English](https://github.com/jobber2955/gbPool/blob/main/README.md)) ([中文](https://github.com/jobber2955/gbPool/blob/main/README_cn.md))
### GBPool, golang BaiPiao(completely free) proxy pool, a free & simple golang proxy pool module, gathering proxies from all kindness free proxy provider.Developed for individual use, but you can apply it anywhere you want.
### Thanks(also if my project is affecting your service, please post an issue, and we will remove your service immediately.)
- [ihuan](https://ip.ihuan.me/)

### What can it do?
- Switch different proxy
- Easy to implement your own provider
- Support manually enable | disable manager
- Support custom process to a proxy(drop | reuse | expired, you can also implement your own process)
- More and more coming soon with your help.

### Install

    1.Simply add this line below to your go.mod
    github.com/jobber2955/gbPool latest
    2.Copy proxy_pool.yaml.example to the root of your project(of course you can also edit your own config file referring to the existing example file. NOTICE: if you choose to use your own config file, remember to edit the init function in /pool/pool.go)
    3.!!!Remember using the correct proxy config!!!
### Example

    // There is a testing.go file at the root, you can use it for testing
    proxyPool := pool.NewProxyPool(10)
	if err := proxyPool.NewManager("ihuan", &public.IhuanConfig{
		Num:         "5",
		Anonymity:   "",
		Type:        "",
		Post:        "",
		Sort:        "1",
		Port:        "",
		KillPort:    "",
		Address:     "中国",
		Key:         "",
		KillAddress: "",
	}); err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Println(len(proxyPool.ProxyChan))
		proxy := <- proxyPool.ProxyChan
		fmt.Println(proxy.Address)
		time.Sleep(time.Second)

		// Emulate reuse proxy, be careful, if the channel is full, this will try forever
		// you need to check the length of channel first
		t := rand.Intn(9)
		if t > 5 && len(proxyPool.ProxyChan) < 10 {
			proxy.ReUse(proxyPool.ProxyChan)
		}
	}

### Things that you need to know
- This is a personally developed module, the main purpose is to help people get a more stable & anonymous way to gather the data they need.
- This module is not for abruptly annoying or even attack those free proxy providers. In contrary, I hope this package can reduce the pressure for them, as this project is more steerable & less request.
- Please be known this is a personal & individual work, some problem may not be answered or solved immediately. Welcome to contribute!
- Using this project is on your own responsibility, any consequence is yours.
- I'm quite a noob in programming, so if you have any bug | suggestion, please let me know in the issue.

### TODO
- ~~Maybe change config file to variable when create new manager.(already changed, and no longer use logger, now we use errors)~~
- Support more provider
- Specify the description of function
