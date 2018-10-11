package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "shippy/consignment-service/proto/consignment"
)

const (
	PORT = ":50051"
)

type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) //存放货物
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}
func (repo *Repository) getAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignment(ctx context.Context, in *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(in)
	if err != nil {
		return nil, err
	}
	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil

}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	allConsignments := s.repo.getAll()
	resp := &pb.Response{Consignments: allConsignments}
	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listen on %s\n", PORT)
	server := grpc.NewServer()
	repo := Repository{}
	pb.RegisterShippingServiceServer(server, &service{repo})
	if err := server.Serve(listener); err != nil {
		fmt.Printf("failed to server:%v\n", err)
	}

}
