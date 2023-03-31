package grpcclient

import (
	"context"
	"fmt"
	"log"
	"os"

	sample_grpc "github.com/Gavachas/grpc_sample/grpc_s"
	"google.golang.org/grpc"
)

const DEFAULT_ADR_GRPC = "localhost:9090"

func GetUserRegion(id int) (string, error) {

	fmt.Println("ClientGRPC")

	portGRPC := os.Getenv("ADR_GRPC")
	if portGRPC == "" {
		portGRPC = DEFAULT_ADR_GRPC
	}

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(portGRPC, opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close() // Maybe this should be in a separate function and the error handled?

	c := sample_grpc.NewItilServiceClient(cc)

	// read Region
	fmt.Println("Reading the region")
	readRegionReq := &sample_grpc.GetUserRequest{Id: int32(id)}
	readRegionRes, readRegionErr := c.GetUserRegion(context.Background(), readRegionReq)
	if readRegionErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readRegionErr)
	}

	fmt.Printf("Region was read: %v \n", readRegionRes)
	return readRegionRes.Name, nil

}
