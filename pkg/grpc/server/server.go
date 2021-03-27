package server

import (
	"context"
	"net"

	"fiy/common/log"
	"fiy/pkg/grpc/proto/host"

	"google.golang.org/grpc"
)

/*
  @Author : lanyulei
*/

const (
	port = ":50051"
)

type server struct {
	host.UnimplementedHostInfoServer
}

func (s *server) GetHostInfo(ctx context.Context, in *host.HostInfoRequest) (*host.HostInfoReply, error) {
	data := in.GetData()
	if data != "" {
		err := insertData(data)
		if err != nil {
			log.Error("插入数据错误，", err)
			return &host.HostInfoReply{Status: false}, err
		}
	}
	return &host.HostInfoReply{Status: true}, nil
}

func RunServer() {
	log.Info("Start the rpc server...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	host.RegisterHostInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
