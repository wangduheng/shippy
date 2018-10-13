package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"
	"os"
	pb "shippy/consignment-service/proto/consignment"
)

const (
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
	fmt.Println("client main1")
	cmd.Init()
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	fmt.Println("client main2")
	// conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()
	// client := pb.NewShippingServiceClient(conn)
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}
	fmt.Println("client main3")
	consignment, err := parseFile(infoFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("client main4")
	resp, err := client.CreateConsignment(context.TODO(), consignment)
	if err != nil {
		panic(err)
	}
	fmt.Println("client main5")
	fmt.Printf("create %t\n", resp.Created)

	resps, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("client main6")
	for _, item := range resps.Consignments {
		fmt.Printf("%v\n", item)
	}
	fmt.Println("client main7")
}
