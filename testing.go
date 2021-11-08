package gbPool
// if you want to use this to test, remember change package to main



//import (
//	"fmt"
//	"github.com/jobber2955/gbPool/pool"
//	"github.com/jobber2955/gbPool/public"
//	"math/rand"
//	"time"
//)
//
//func main() {
//
//	proxyPool := pool.NewProxyPool(10)
//	if err := proxyPool.NewManager("ihuan", &public.IhuanConfig{
//		Num:         "5",
//		Anonymity:   "",
//		Type:        "",
//		Post:        "",
//		Sort:        "1",
//		Port:        "",
//		KillPort:    "",
//		Address:     "中国",
//		Key:         "",
//		KillAddress: "",
//	}); err != nil {
//		fmt.Println(err)
//		return
//	}
//	for {
//		fmt.Println(len(proxyPool.ProxyChan))
//		proxy := <- proxyPool.ProxyChan
//		fmt.Println(proxy.Address)
//		time.Sleep(time.Second)
//
//		// Emulate reuse proxy, be careful, if the channel is full, this will try forever
//		// you need to check the length of channel first
//		t := rand.Intn(9)
//		if t > 5 && len(proxyPool.ProxyChan) < 10 {
//			proxy.ReUse(proxyPool.ProxyChan)
//		}
//	}
//}
