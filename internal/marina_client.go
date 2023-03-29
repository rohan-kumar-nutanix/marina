/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * GRPC client for echo service.
 */

package internal

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marina "github.com/nutanix-core/content-management-marina/protos/marina"
	glog "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func init() {
	flag.Set("logtostderr", "true")
}

func marinaClient(port uint) (*grpc.ClientConn, marina.MarinaClient) {
	//address := fmt.Sprintf("127.0.0.1:%d", port)
	// address := fmt.Sprintf("10.37.161.66:%d", port)
	pcIp := "localhost" // "10.96.16.100" // "0.0.0.0" //"10.33.33.78"
	port = 9200
	address := fmt.Sprintf("%s:%d", pcIp, port)
	fmt.Println("Connecting to GRPC server at ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Failed grpc server connection: %v.", err)
		glog.Fatalf("Failed grpc server connection: %v.", err)
	}

	// Create new client to talk to Echo server.
	mclient := marina.NewMarinaClient(conn)
	return conn, mclient
}

func testMarina(client marina.MarinaClient, ctx context.Context) {
	// arg := &marina.CatalogItemGetArg{}

	response, err := client.CatalogItemGet(ctx, &marina.CatalogItemGetArg{})
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	// fmt.Printf("Response	 received \n Catalog Items :%s  \n Description", response)
	for i, item := range response.CatalogItemList {
		fmt.Printf("Catalog Item %v: %v\n", i, item)
	}

}

func GetAllCatalogItems(client marina.MarinaClient, ctx context.Context) {

	response, err := client.CatalogItemGet(ctx, &marina.CatalogItemGetArg{})
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	// fmt.Printf("Response	 received \n Catalog Items :%s  \n Description", response)
	for i, item := range response.CatalogItemList {
		// fmt.Printf("Catalog Item %v: %v\n", i, item)
		fmt.Printf("%v: %v\n", i, item)
	}

}

func GetCatalogItemsById(client marina.MarinaClient, ctx context.Context, uuids []string) {
	guuidList := []string{"ae08ab22-08ac-4da4-9c5f-f4d3d8e2f6a4", "d9d3a431-8dbc-483f-9ce3-b3ce23be84a7",
		"e29784c7-8b45-46d2-8137-a45d6143e00e"}
	/* guuid1, _ := uuid4.StringToUuid4("e04a5278-f9ed-41cb-8601-89ab60c7f75f")
	item1 := &marina.CatalogItemId{
		GlobalCatalogItemUuid: guuid1.RawBytes(),
	} */

	/* invalid_item1 := &marina.CatalogItemId{
		GlobalCatalogItemUuid: []byte("12111"),
	} */
	var items []*marina.CatalogItemId
	// items = append(items, item1) //, invalid_item1)
	for _, gcid := range guuidList {
		guuid, _ := uuid4.StringToUuid4(gcid)
		ver := new(int64)
		*ver = 1
		item := &marina.CatalogItemId{
			GlobalCatalogItemUuid: guuid.RawBytes(),
			// Version:               ver,
		}
		items = append(items, item)
	}
	// items = append(items, invalid_item1)
	var itemTypes []marina.CatalogItemInfo_CatalogItemType
	itemTypes = append(itemTypes, *marina.CatalogItemInfo_kVmTemplate.Enum())

	arg := &marina.CatalogItemGetArg{
		CatalogItemIdList:   items,
		CatalogItemTypeList: itemTypes,
	}
	// response, err := client.CatalogItemGet(ctx, &marina.CatalogItemGetArg{})
	fmt.Println("Fetching CatalogItems by UUIDs")
	response, err := client.CatalogItemGet(ctx, arg)
	// fmt.Printf("Error occured %v", err)
	if err != nil {
		fmt.Printf("Marina request error: %s\n", err)
		errStatus, _ := status.FromError(err)
		fmt.Printf("Error Status Message : %s \nCode %s \n full obj : %v", errStatus.Message(),
			errStatus.Code(), errStatus)

		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("********Catalog Items********\n-------------------------\n")
	for i, item := range response.CatalogItemList {
		fmt.Printf("'%v': %v\n", i, item)
	}

}

func main1() {

	// TODO: Consider moving these to init().
	var grpcServerPort uint
	grpcServerPort = 9200 //32391 //9200 //30188 //9200
	//glog.Infof("gRPC server port to connect: %v", grpcServerPort)
	// fmt.Println("gRPC server port to connect: ", grpcServerPort)
	conn, mclient := marinaClient(grpcServerPort)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// fmt.Println("Making echo server request.")
	// testEcho(client, ctx)

	fmt.Println("Making Marina server request CatalogItemGet.")
	GetAllCatalogItems(mclient, ctx)
	fmt.Println("-----------------------------------------------")
	fmt.Println("Fetching CatalogItemGet by UUID's:")
	fmt.Println("-----------------------------------------------")
	GetCatalogItemsById(mclient, ctx, nil)
}
