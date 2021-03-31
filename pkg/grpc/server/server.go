package server

import (
	"context"
	"net"

	"fiy/common/log"
	pb "fiy/pkg/grpc/proto"

	"google.golang.org/grpc"
)

/*
  @Author : lanyulei
*/

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedHostInfoServer
}

func (s *server) GetHostInfo(ctx context.Context, in *pb.HostInfoRequest) (*pb.CommonReply, error) {
	data := in.GetData()
	if data != "" {
		err := insertData(data)
		if err != nil {
			log.Error("插入数据错误，", err)
			return &pb.CommonReply{Status: false}, err
		}
	}
	return &pb.CommonReply{Status: true}, nil
}

func RunServer() {
	log.Info("Start the rpc server...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHostInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
