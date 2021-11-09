# GBPool-- 简单好用的Go免费代理池
## 简介([English](https://github.com/jobber2955/gbPool/blob/main/README.md)) ([中文](https://github.com/jobber2955/gbPool/blob/main/README_cn.md))
### GBPool, Golang baipiao proxy pool,一个免费、简单的代理池模块,它从免费代理提供商处获取代理.主要用作个人使用, 但你想怎么用就怎么用.
### 感谢(同时如果本模块对您的服务造成了影响，请及时提出issue，我会第一时间删除)
- [iHuan](https://ip.ihuan.me/)

### 它能做什么?
- 实时切换代理池
- 十分简单地实现其他供应商
- 支持启用、禁用manager
- 支持对代理IP的操作(删除、复用、过期丢弃，你也可以自行实现其他操作)
- 在你的帮助下，还会有更多的功能！

### 安装

    1.将以下命令添加到go.mod文件即可
    github.com/jobber2955/gbPool latest
    2.把proxy_pool.yaml.example文件复制到你项目的根目录下，当然你也可以将example文件中的内容集成到你自己的配置文件中。注意：如果选择集成的方式，别忘了修改/pool/pool.go文件中的配置读取设置
    3.记得使用正确的Config！
### 样例(每种Fetcher的Config配置在下方)

    import (
        "fmt"
        "github.com/jobber2955/gbPool/pool"
        "github.com/jobber2955/gbPool/public"
        "time"
    )

    func main() {
        proxyPool := pool.NewProxyPool(30)
        err := proxyPool.NewManager("ihuan", &public.IHuanConfig{
        Num:         "30",
        Anonymity:   "",
        Type:        "",
        Post:        "",
        Sort:        "",
        Port:        "",
        KillPort:    "",
        Address:     "",
        Key:         "",
        KillAddress: "",
        })
        if err != nil {
            return
        }
        for {
            fmt.Println(<- proxyPool.ProxyChan)
            time.Sleep(time.Second)
        }
    }

### Config
- iHuan

<table>
    <tr>
        <th>属性名</th>
        <th>变量类型</th>
        <th>含义</th>
        <th>示例</th>
        <th>备注</th>
    </tr>
    <tr>
        <td>Num</td>
        <td>string</td>
        <td>每次获取的代理数量</td>
        <td>10</td>
        <td rowspan="10">iHuan使用表单做POST，而非RESTful的JSON，会有空值(跟Golang的nil一样)，所以用string会很方便，如果用类似int的话，默认是0而不是nil</td>
    </tr>
    <tr>
        <td>Anonymity</td>
        <td>string</td>
        <td>匿名程度</td>
        <td>0:透明代理<br>1:普通代理<br>2:高匿代理<br>留空: 全部</td>
    </tr>
    <tr>
        <td>Type</td>
        <td>string</td>
        <td>支持HTTPS || HTTP</td>
        <td>0:只允许HTTP<br>1:只允许HTTPS<br>留空: 全部</td>
    </tr>
    <tr>
        <td>Post</td>
        <td>string</td>
        <td>是否支持POST</td>
        <td>1:是<br>留空: 全部</td>
    </tr>
    <tr>
        <td>Sort</td>
        <td>string</td>
        <td>排序方式</td>
        <td>1:验证时间从近到远<br>2:验证时间从远到近<br>3:存活时间由短到长<br>4:存活时间由长到短<br>留空: 随机</td>
    </tr>
    <tr>
        <td>Port</td>
        <td>string</td>
        <td>指定端口</td>
        <td>8080(只支持一个)</td>
    </tr>
    <tr>
        <td>KillPort</td>
        <td>string</td>
        <td>指定排除端口</td>
        <td>8080(只支持一个)</td>
    </tr>
    <tr>
        <td>Address</td>
        <td>string</td>
        <td>指定地址</td>
        <td>中国/美国(只支持一个，不确定非简体中文是否可以)</td>
    </tr>
    <tr>
        <td>排除指定地址</td>
        <td>string</td>
        <td>Specific exclude Address</td>
        <td>中国/美国(只支持一个，不确定非简体中文是否可以)</td>
    </tr>
    <tr>
        <td>Key</td>
        <td>string</td>
        <td>iHuan POST验证Key</td>
        <td>无需填写</td>
    </tr>
</table>

### 注意事项
- 这是个个人开发的模块，主要目的是帮助人们更方便、更匿名地获取所需数据。
- 这个模块并不是为了暴力爬取，相反，由于本模块更可控、更文明，反而可能会帮助服务商减轻被爬压力。
- 由于个人精力有限，可能有些问题无法及时回答或没有能力解决。
- 使用本项目均为你的个人意愿，一切造成的后果由你自行承担。
- 我是编程小白，有任何程序、设计或其他方面的bug、建议，欢迎在issue中提出。

### TODO
- ~~也许将配置文件改为参数传入(已完成，并不再使用logger，取而代之的是返回error)~~
- 支持更多提供商
- 为函数的用途补充说明
- 也许我应该将公共结构体放到根目录？