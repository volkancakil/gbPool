package gbPool

//import (
//	"fmt"
//	"gbPool/pool"
//	"github.com/sirupsen/logrus"
//	"math/rand"
//	"time"
//)
//
//func main() {
//	logger := logrus.New()
//	logger.SetReportCaller(true)
//	proxyPool := pool.NewProxyPool(logger)
//	proxyPool.NewManager("ihuan")
//	for {
//		proxy := <- proxyPool.ProxyChan
//		fmt.Println(proxy.Address)
//		time.Sleep(time.Second)
//		t := rand.Intn(9)
//		if t > 5 {
//			proxy.ReUse(proxyPool.ProxyChan)
//		}
//	}
//}
