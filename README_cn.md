# GBPool-- 简单好用的Go免费代理池
## 简介([English](https://github.com/jobber2955/gbPool/blob/main/README.md)) ([中文](https://github.com/jobber2955/gbPool/blob/main/README_cn.md))
### GBPool, Golang baipiao proxy pool,一个免费、简单的代理池模块,它从免费代理提供商处获取代理.主要用作个人使用, 但你想怎么用就怎么用.
### 感谢(同时如果本模块对您的服务造成了影响，请及时提出issue，我会第一时间删除)
- [ihuan](https://ip.ihuan.me/)

### 它能做什么?
- 实时切换代理池
- 十分简单地实现其他供应商
- 支持启用、禁用manager
- 支持对代理IP的操作(删除、复用、过期丢弃，你也可以自行实现其他操作)
- 在你的帮助下，还会有更多的功能！

### 安装

    // 1.将以下命令添加到go.mod文件即可
    github.com/jobber2955/gbPool latest
    // 2.把proxy_pool.yaml.example文件复制到你项目的根目录下，当然你也可以将example文件中的内容集成到你自己的配置文件中。注意：如果选择集成的方式，别忘了修改/pool/pool.go文件中的配置读取设置
### 样例

    // 模块根目录下有哥testing.go文件可以拿来测试。
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

### 注意事项
- 这是个个人开发的模块，主要目的是帮助人们更方便、更匿名地获取所需数据。
- 这个模块并不是为了暴力爬取，相反，由于本模块更可控、更文明，反而可能会帮助服务商减轻被爬压力。
- 由于个人精力有限，可能有些问题无法及时回答或没有能力解决。
- 使用本项目均为你的个人意愿，一切造成的后果由你自行承担。
- 我是编程小白，有任何程序、设计或其他方面的bug、建议，欢迎在issue中提出。

### TODO
- 也许将配置文件改为参数传入
- 支持更多提供商
- 为函数的用途补充说明
