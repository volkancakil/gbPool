# GBPool-- a simple but useful golang free proxy pool
## Intro([English](https://github.com/jobber2955/gbPool/blob/main/README.md)) ([中文](https://github.com/jobber2955/gbPool/blob/main/README_cn.md))
### GBPool, golang BaiPiao(completely free) proxy pool, a free & simple golang proxy pool module, gathering proxies from all kindness free proxy provider.Developed for individual use, but you can apply it anywhere you want.
### Thanks(also if my project is affecting your service, please post an issue, and we will remove your service immediately.)
- [iHuan](https://ip.ihuan.me/)

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

### Example(Config for each type of manager's fetcher is down below)

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

### Config Type
- iHuan

<table>
    <tr>
        <th>Name</th>
        <th>Type</th>
        <th>Meaning</th>
        <th>Value Example</th>
        <th>Comment</th>
    </tr>
    <tr>
        <td>Num</td>
        <td>string</td>
        <td>number of proxy to get each time</td>
        <td>10</td>
        <td rowspan="10">iHuan use form to post rather than RESTful API. Needs to have an empty value of each key(same as nil, but for types like int, their default "ZERO VALUE" is 0, so just use string for simple)</td>
    </tr>
    <tr>
        <td>Anonymity</td>
        <td>string</td>
        <td>Level of the proxy's anonymity</td>
        <td>0:Transparent proxy<br>1:Normal proxy<br>2:High anonymity proxy<br>empty: All</td>
    </tr>
    <tr>
        <td>Type</td>
        <td>string</td>
        <td>HTTPS || HTTP</td>
        <td>0:HTTP only<br>1:HTTPS only<br>empty: All</td>
    </tr>
    <tr>
        <td>Post</td>
        <td>string</td>
        <td>Support Post method</td>
        <td>1:Support<br>empty: All</td>
    </tr>
    <tr>
        <td>Sort</td>
        <td>string</td>
        <td>Sorting Type</td>
        <td>1:Validate time nearest first<br>2:Validate time farthest first<br>3:Alive duration shortest first<br>4:Alive duration longest first<br>empty: Random</td>
    </tr>
    <tr>
        <td>Port</td>
        <td>string</td>
        <td>Specific port</td>
        <td>8080(only 1 port supported)</td>
    </tr>
    <tr>
        <td>KillPort</td>
        <td>string</td>
        <td>Specific exclude port</td>
        <td>8080(only 1 port supported)</td>
    </tr>
    <tr>
        <td>Address</td>
        <td>string</td>
        <td>Specific Address</td>
        <td>中国/美国(only one support, not sure about non-simplified-chinese words)</td>
    </tr>
    <tr>
        <td>KillAddress</td>
        <td>string</td>
        <td>Specific exclude Address</td>
        <td>中国/美国(only one support, not sure about non-simplified-chinese words)</td>
    </tr>
    <tr>
        <td>Key</td>
        <td>string</td>
        <td>iHuan post validation key</td>
        <td>You don't have to specify this </td>
    </tr>
</table>

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
- Maybe I should change all the public struct(like &IhuanConfig) to the root directory?
