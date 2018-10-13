package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	pb "shippy/consignment-service/proto/consignment"
	vesselProto "shippy/vessel-service/proto/vessel"
)

type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) //存放货物
	GetAll() []*pb.Consignment                                   //获取长湖中所有的货物
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	// 检查是否有适合的货轮
	vReq := &vesselProto.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}
	fmt.Printf("found vessel:%s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = consignment
	return nil

}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	resp.Consignments = allConsignments
	return nil
}

func main() {

	server := micro.NewService(micro.Name("go.micro.srv.consignment"), micro.Version("latest"))
	server.Init()
	repo := Repository{}
	vClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", server.Client())

	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})
	if err := server.Run(); err != nil {
		panic(err)
	}

}
