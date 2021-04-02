package client

import (
	"context"
	"fiy/common/log"
	pb "fiy/pkg/grpc/proto"
	"time"

	"google.golang.org/grpc"
)

/*
  @Author : lanyulei
*/

func (r *rpcClient) RunClient(address string, interval int) {
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
			c := pb.NewHostInfoClient(conn)

			// Contact the server and print out its response.
			log.Info("开始收集资源数据...")
			data, err := r.IntegrateData()
			if err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			log.Info("开始上报资源数据...")
			_, err = c.GetHostInfo(ctx, &pb.HostInfoRequest{Data: data})
			if err != nil {
				log.Fatalf("开始上报资源数据: %v", err)
				return
			}
		}()

		t.Reset(td)
	}
}
