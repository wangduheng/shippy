package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	// "log"
	"os"
	pb "shippy/consignment-service/proto/consignment"
)

const (
	ADDRESS           = "149.28.17.62:50051"
	DEFAULT_INFO_FILE = "consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)

	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {

	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewShippingServiceClient(conn)
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}
	consignment, err := parseFile(infoFile)
	if err != nil {
		panic(err)
	}
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("create %t\n", resp.Created)

	resps, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		panic(err)
	}
	for _, item := range resps.Consignments {
		fmt.Printf("%v\n", item)
	}
}
