# GBPool-- a simple but useful golang free proxy pool
## Intro([English](https://github.com/jobber2955/gbPool/blob/main/README.md)) ([中文](https://github.com/jobber2955/gbPool/blob/main/README_cn.md))
### GBPool, golang baipiao proxy pool, a free & simple golang proxy pool module, gathering proxies from those kindness free proxy provider.Developed for indiviual use, but you can apply it anywhere you want.
### Thanks(also if my project is affecting your service, please post a issue and i will remove your service immediately.)
- [ihuan](https://ip.ihuan.me/)

### What can it do?
- Switch different proxy
- Easy to implement your own provider
- Support manually enable | disable manager
- Support custom process to a proxy(drop | reuse | expired, you can also implement your own process)
- More and more coming soon with your help.

### Install

    // 1.Simply add this line below to your go.mod
    github.com/jobber2955/gbPool latest
    // 2.Copy proxy_pool.yaml.example to the root of your project(of course you can also edit your own config file referring to the existing example file. NOTICE: if you choose to use your own config file, remember to edit the init function in /pool/pool.go)

### Example

    // There is a testing.go file at the root, you can use it for testing
    logger := logrus.New()
	logger.SetReportCaller(true)
	proxyPool := pool.NewProxyPool(logger)
	proxyPool.NewManager("ihuan")
	for {
		proxy := <- proxyPool.ProxyChan
		fmt.Println(proxy.Address)
		time.Sleep(time.Second)
		t := rand.Intn(9)
		if t > 5 {
			proxy.ReUse(proxyPool.ProxyChan)
		}
	}

### Things that you need to know
- This is a personally developed module, the main purpose is to help people get a more stable & anonymous way to gather the data they need.
- This module is not for abruptly annoying or even attack those free proxy providers. In contrary, i hope this package can reduce the pressure for them, as this project is more steerable & less request.
- Please be known this is a personal & individual work, some problem may not be answered or solved immediately. Welcome to contribute!
- Using this project is on your own responsibility, any consequence is yours.
- I'm quite a noob in programming, so if you have bug | suggestion, please let me know in the issue.

### TODO
- Maybe change config file to variable when create new manager.
- Support more provider
- Specify the description of function
