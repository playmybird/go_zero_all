package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-examples/rpc/remote/unary"
)

const timeFormat = "15:04:05"

var configFile = flag.String("f", "config.yaml", "the config file")

func main() {
	flag.Parse()

	var c zrpc.RpcClientConf
	conf.MustLoad(*configFile, &c)

	client := zrpc.MustNewClient(c)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn := client.Conn()
			greet := unary.NewGreeterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			resp, err := greet.Greet(ctx, &unary.Request{
				Name: "kevin",
			})
			if err != nil {
				fmt.Printf("%s X %s\n", time.Now().Format(timeFormat), err.Error())
			} else {
				fmt.Printf("%s => %s\n", time.Now().Format(timeFormat), resp.Greet)
			}
			cancel()
		}
	}
}
