package client

import (
	"context"
	"fiy/common/log"
	"fiy/pkg/grpc/proto/host"
	"time"

	"google.golang.org/grpc"
)

/*
  @Author : lanyulei
*/

const (
	address = "localhost:50051"
)

func RunClient(address string, interval int) {
	// 定时收集上报数据，每隔5分钟上报一次数据
	td := time.Duration(interval) * time.Minute
	t := time.NewTicker(td)
	defer t.Stop()

	log.Info("start...")
	for {
		<-t.C

		go func() {
			// Set up a connection to the server.
			log.Info("创建服务端连接...")
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
				return
			}
			defer conn.Close()
			c := host.NewHostInfoClient(conn)

			// Contact the server and print out its response.
			log.Info("开始收集资源数据...")
			data, err := IntegrateData()
			if err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			log.Info("开始上报资源数据...")
			_, err = c.GetHostInfo(ctx, &host.HostInfoRequest{Data: data})
			if err != nil {
				log.Fatalf("收集资源数据失败: %v", err)
				return
			}
		}()

		t.Reset(td)
	}
}
