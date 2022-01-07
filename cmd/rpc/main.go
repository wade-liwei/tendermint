package main

import (
	"context"
	"fmt"

	"flag"
	"sync"
	"time"

	"crypto/md5"
	"sync/atomic"

	tmrand "github.com/tendermint/tendermint/libs/rand"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"golang.org/x/time/rate"
)

func getHTTPClient(rpcAddr string) *rpchttp.HTTP {
	c, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		panic(err)
	}
	return c
}

func MakeTxKV() ([]byte, []byte, []byte) {
	k := []byte(tmrand.Str(8))
	v := []byte(tmrand.Str(8))
	return k, v, append(k, append([]byte("="), v...)...)
}

func main() {

	host := flag.String("host", "http://127.0.0.1", "请输入host地址")
	port := flag.Int("port", 26657, "请输入端口号")

	md5Count := flag.Int("md5count", 5000, "请输入md5个数")
	sendRate := flag.Int("sendRate", 1000, "单位是Microsecond = 1 millisecond = 0.1second ")
	duration := flag.Int("duration", 10, "单位是 second ")

	flag.Parse()

	md5hashs := make([]byte, 0, md5.Size*(*md5Count))
	for i := 0; i < *md5Count; i++ {
		hash := md5.Sum([]byte(fmt.Sprintf("%d", i)))
		md5hashs = append(md5hashs, hash[:]...)
	}

	c := getHTTPClient(fmt.Sprintf("%s:%d", *host, *port))

	limiter := rate.NewLimiter(rate.Every(time.Duration(1000*1000 / *sendRate)*time.Microsecond), 5)

	wg := sync.WaitGroup{}
	maxChannel := make(chan struct{}, 10000)
	beginTime := time.Now()
	cxt, _ := context.WithCancel(context.TODO())
	var count int64

	for {
		limiter.Wait(cxt)
		wg.Add(1)

		go func() {
			maxChannel <- struct{}{}
			bres, err := c.BroadcastTxSync(context.Background(), append([]byte(tmrand.Str(10)+"="), md5hashs...))
			if err != nil {
				panic(err)
			}
			_ = bres
			wg.Done()
			<-maxChannel
			count = atomic.AddInt64(&count, 1)
		}()

		if beginTime.Add(time.Duration(*duration)*time.Second).Before(time.Now()) && count >= int64((*sendRate)*(*duration)) {
			break
		}
	}

	wg.Wait()

	fmt.Printf("every: %v micros duration: %v txcount: %v real duration: %v  md5hashs.len: %vB/tx  %vKB/tx totalSize %v KB / realDuration %v = %v KB/s  \n", 1000*1000 / *sendRate, *duration, count, time.Now().Sub(beginTime), len(md5hashs), len(md5hashs)/1024, int64(len(md5hashs)/1024)*count, int64(time.Now().Sub(beginTime)/time.Second), int64(len(md5hashs)/1024)*count/int64(time.Now().Sub(beginTime)/time.Second))

	fmt.Printf("total size: %vKB  %vMB  %vGB  \n", int64(len(md5hashs)/1024)*count, int64(len(md5hashs)/1024)*count/1024, int64(len(md5hashs)/1024)*count/1024/1024)
}
