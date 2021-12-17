package main

import (
	"context"
	"fmt"

	"encoding/json"
	"flag"

	//"crypto/md5"

	"crypto/md5"

	tmrand "github.com/tendermint/tendermint/libs/rand"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
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
	randomPre := flag.Bool("random", false, "默认不启用，启用后增加随机的xxxxxxxx=")
	md5Count := flag.Int("md5count", 0, "请输入md5个数")
	txCount := flag.Int("txcount", 0, "请输入md5个数")
	flag.Parse()

	if *md5Count == 0 || *txCount == 0 {
		fmt.Printf("please  with md5 %d  and tx  %d  flag  port %d  \n",*md5Count,*txCount,*port)
		return
	}

	md5hashs := make([]byte, 0, md5.Size*(*md5Count))
	for i := 0; i < *md5Count; i++ {
		hash := md5.Sum([]byte(fmt.Sprintf("%d", i)))
		md5hashs = append(md5hashs, hash[:]...)
	}

	c := getHTTPClient(fmt.Sprintf("%s:%d", *host, *port))
	if *randomPre {
		for i := 0; i < *txCount; i++ {
			bres, err := c.BroadcastTxSync(context.Background(), append([]byte(tmrand.Str(10)+"="), md5hashs...))

			if err != nil {
				panic(err)
			}

			b, err := json.Marshal(bres)
			if err != nil {
				fmt.Println("JSON ERR:", err)
			}
			fmt.Println(string(b))
		}

	} else {

		for i := 0; i < *txCount; i++ {
			c := getHTTPClient(fmt.Sprintf("%s:%d", *host, *port))
			bres, err := c.BroadcastTxSync(context.Background(), md5hashs)

			if err != nil {
				panic(err)
			}

			b, err := json.Marshal(bres)
			if err != nil {
				fmt.Println("JSON ERR:", err)
			}
			fmt.Println(string(b))
		}
	}
}

