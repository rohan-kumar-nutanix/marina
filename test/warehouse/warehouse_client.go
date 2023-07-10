/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * GRPC client for echo service.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	glog "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	warehouse "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
)

func init() {
	flag.Set("logtostderr", "true")
}

// var pcip = "10.97.188.105" // "localhost" // "10.97.16.174"
var pcip = "localhost"

func warehouseClient(port uint) (*grpc.ClientConn, warehouse.WarehouseServiceClient) {
	// address := fmt.Sprintf("127.0.0.1:%d", port)
	// address := fmt.Sprintf("10.37.161.66:%d", port)
	// pcIp := "localhost" // pcip // //"10.96.16.100" // "0.0.0.0" //"10.33.33.78"
	// pcIp := "10.97.188.105"
	pcIp := "localhost"
	port = 9200
	address := fmt.Sprintf("%s:%d", pcIp, port)
	fmt.Println("Connecting to GRPC server at ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Failed grpc server connection: ", err)
		glog.Fatalf("Failed grpc server connection: %v.", err)
	}

	// Create new client to talk to Echo server.
	mclient := warehouse.NewWarehouseServiceClient(conn)
	// marina.NewMarinaClient(conn)
	return conn, mclient
}

func warehouseItemClient(port uint) (*grpc.ClientConn, warehouse.WarehouseItemsServiceClient) {
	// address := fmt.Sprintf("127.0.0.1:%d", port)
	// address := fmt.Sprintf("10.37.161.66:%d", port)
	// pcIp := "10.97.188.105" // "10.96.16.100" // "0.0.0.0" //"10.33.33.78"
	pcIp := "localhost"
	// port = 9200
	address := fmt.Sprintf("%s:%d", pcIp, port)
	fmt.Println("Connecting to GRPC server at ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Failed grpc server connection: ", err)
		glog.Fatalf("Failed grpc server connection: %v.", err)
	}

	// Create new client to talk to Echo server.
	mclient := warehouse.NewWarehouseItemsServiceClient(conn)
	// marina.NewMarinaClient(conn)
	return conn, mclient
}

func testGetWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	response, err := client.GetWarehouse(ctx, &warehouse.GetWarehouseArg{
		ExtId: proto.String("87e24150-dc7c-46e8-568d-2de3750d919d")})
	if err != nil {
		fmt.Println("Marina request error:", err)
		s, ok := status.FromError(err)
		fmt.Printf("Marina request error s=%v ok=%v ", s, ok)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Warehouse Response ", response)
	fmt.Println("Warehouse Response details", response)

}

func testDeleteWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	arg := &warehouse.DeleteWarehouseArg{ExtId: proto.String("7db350a6-3e18-4723-5a0f-a5c076bc17b2")}
	response, err := client.DeleteWarehouse(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		s, ok := status.FromError(err)
		fmt.Printf("Marina request error s=%v ok=%v ", s, ok)
		return
	}
	fmt.Printf("Warehouse Delete arg: %v\n", arg)
	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Warehouse Delete Response ", response)
	fmt.Println("Warehouse Delete Response details", response)

}

func deleteAllWarehouses(ctx context.Context, client warehouse.WarehouseServiceClient) {

	ret := testListWarehouse(ctx, client)
	resp := ret.Content.GetWarehouseArrayData().GetValue()

	for _, w := range resp {
		fmt.Printf("Deleting Warehouse with UUID %s", *w.Name)
		// arg := &warehouse.DeleteWarehouseArg{ExtId: proto.String("5e81dd9e-fe69-4323-58e4-455dc9f8d240")}
		if w.GetBase() != nil && w.GetBase().ExtId != nil {
			fmt.Printf("Warehouse Ext ID : %s", *w.Base.ExtId)
			arg := &warehouse.DeleteWarehouseArg{ExtId: proto.String(*w.Base.ExtId)}
			response, err := client.DeleteWarehouse(ctx, arg)
			if err != nil {
				fmt.Println("Marina request error:", err)
				s, ok := status.FromError(err)
				fmt.Printf("Marina request error s=%v ok=%v ", s, ok)
				return
			}
			fmt.Printf("Warehouse Delete arg: %v\n", arg)
			fmt.Println("***********Response received from server*******")
			fmt.Println("-----------------------------------------------")
			fmt.Println("Warehouse Delete Response ", response)
			fmt.Println("Warehouse Delete Response details", response)
		}
		fmt.Println("This Warehouse Don't have ExtID")

	}

}

func testListWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) *warehouse.ListWarehousesRet {
	response, err := client.ListWarehouses(ctx, &warehouse.ListWarehousesArg{})
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return nil
	}
	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Warehouse List Response ", response)
	return response

}

func testListWarehouseItem(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {
	arg := &warehouse.ListWarehouseItemsArg{
		ExtId: proto.String("4e98f7ea-3816-49c7-46f7-46d2e7e6a0ae"), // 3989cf51-d046-4ee9-77e1-4dc9339927a8"),
	}

	glog.Infof("WarehouseItem Arg %v", arg)
	response, err := client.ListWarehouseItems(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("WarehouseItem List Response ", response)

}

func testGetWarehouseItemById(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {
	arg := &warehouse.GetWarehouseItemByIdArg{
		ExtId:          proto.String("4e98f7ea-3816-49c7-46f7-46d2e7e6a0ae"),
		WarehouseExtId: proto.String("74865c8a-c99c-4ffb-588d-1843c2f5dbbe"),
	}

	glog.Infof("WarehouseItemById Arg %v", arg)
	response, err := client.GetWarehouseItemById(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("WarehouseItemGetById Response ", response)

}
func testCreateWarehouseItem(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {
	// sourceRef1 := &warehouse.WarehouseItem_CloudObjectSourceSourceReference {CloudObjectSourceSourceReference: warehouse.CloudObjectSource{ObjectKey: }}

	// sourceRef := &warehouse.CloudObjectSourceWrapper{
	// 	Value: &warehouse.CloudObjectSource{ObjectKey: proto.String("Alpine_Docker_Image")},
	// }

	sourceRef := &warehouse.WarehouseItem_DockerImageReferenceSourceReference{
		DockerImageReferenceSourceReference: &warehouse.DockerImageReferenceWrapper{
			Value: &warehouse.DockerImageReference{
				Sha256HexDigest: proto.String("sha256:0ef2e08ed3fabfc44002ccb846c4f2416a2135affc3ce39538834059606f32dd"),
			},
		},
	}

	/*arg := &warehouse.AddItemToWarehouseArg{
		ExtId: proto.String("4e98f7ea-3816-49c7-46f7-46d2e7e6a0ae"),
		Body: &warehouse.WarehouseItem{
			Base:        nil,
			Name:        proto.String("Ngnix Docker Image "),
			Description: proto.String("latest version fo Ngnix Docker Image"),
			Type:        warehouse.ItemTypeMessage_DOCKER_IMAGE.Enum(),
			SizeBytes:   proto.Int64(100),
			MetaData:    nil,
		},
	}*/

	arg := &warehouse.AddItemToWarehouseArg{
		ExtId: proto.String("7db350a6-3e18-4723-5a0f-a5c076bc17b2"),
		Body: &warehouse.WarehouseItem{
			Name:            proto.String("Ngnix-Docker-Image"),
			Description:     proto.String("Latest Nginx Docker"),
			Type:            warehouse.ItemTypeMessage_CONTAINER_IMAGE.Enum(),
			SizeBytes:       proto.Int64(100),
			SourceReference: sourceRef,
			Signature:       nil,
		},
	}

	glog.Infof("WarehouseItem Arg %v", arg)
	response, err := client.AddItemToWarehouse(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Warehouse Response %v", response)

}

func testDeleteWarehouseItemById(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {
	arg := &warehouse.DeleteWarehouseItemArg{
		ExtId:          proto.String("00f02076-ff62-4a43-5565-7dda0bcaf771"),
		WarehouseExtId: proto.String("7db350a6-3e18-4723-5a0f-a5c076bc17b2"),
	}
	response, err := client.DeleteWarehouseItem(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		s, ok := status.FromError(err)
		fmt.Printf("Marina request error s=%v ok=%v ", s, ok)
		return
	}
	fmt.Printf("WarehouseItem Delete arg: %v\n", arg)
	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("WarehouseItem Delete Response ", response)
	fmt.Println("WarehouseItem Delete Response details", response)

}

func testDeleteAllWarehouseItems(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {
	arg := &warehouse.DeleteWarehouseItemArg{
		ExtId:          proto.String("74865c8a-c99c-4ffb-588d-1843c2f5dbbe"),
		WarehouseExtId: proto.String("4e98f7ea-3816-49c7-46f7-46d2e7e6a0ae"),
	}
	response, err := client.DeleteWarehouseItem(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		s, ok := status.FromError(err)
		fmt.Printf("Marina request error s=%v ok=%v ", s, ok)
		return
	}
	fmt.Printf("WarehouseItem Delete arg: %v\n", arg)
	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("WarehouseItem Delete Response ", response)
	fmt.Println("WarehouseItem Delete Response details", response)

}

func testUpdateWarehouseItemById(ctx context.Context, client warehouse.WarehouseItemsServiceClient) {

	arg := &warehouse.UpdateWarehouseItemMetadataArg{
		ExtId:          proto.String("93f47699-0a93-4e44-636a-0e77c4fb603b"),
		WarehouseExtId: proto.String("9405e481-f27c-4332-55eb-dc15979079fe"),
		Body: &warehouse.WarehouseItem{
			Name:        proto.String("Ocean Container1"),
			Description: proto.String("Updated Container desc"),
		},
	}

	glog.Infof("WarehouseItemUpudate Arg %v", arg)
	response, err := client.UpdateWarehouseItemMetadata(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("WarehouseItemUpdate Response ", response)

}

/*func testCreateWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	config := &warehouse.StorageBacking_CloudStorageConfig{}
	storageWrapper := &warehouse.CloudStorageWrapper{Value: &warehouse.CloudStorage{
		ApiKey:           proto.String("tW_KDPCzz12MZiokdYRULLXlAaVh9ZuL"),
		SecretKey:        proto.String("OUe2waPYFv2QrP9b2tbVSgc04WkwMrIp"),
		AvailabilityZone: proto.String("us-east-1"),
		BucketName:       proto.String("ncr-warehouse-docker-production-images"),
		Type:             warehouse.CloudStorageServiceTypeMessage_NUTANIX.Enum(),
	}}
	config.CloudStorageConfig = storageWrapper
	arg := &warehouse.CreateWarehouseArg{
		Body: &warehouse.Warehouse{
			Name:              proto.String("Broadcom Warehouse II for Test Bed"),
			Description:       proto.String("Broadcom dockers."),
			Metadata:          nil,
			AdditonalMetadata: nil,
			Type:              warehouse.WarehouseTypeMessage_FULL_CONTROL.Enum(),
			StorageBacking: &warehouse.StorageBacking{
				Config: config,
			},
			SecurityPolicies: nil,
			ExtraInfo:        nil,
		},
	}
	glog.Infof("Warehouse Arg %v", arg)
	response, err := client.CreateWarehouse(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	// fmt.Printf("Response	 received \n Catalog Items :%s  \n Description", response)
	fmt.Println("Warehouse Response ", response)

}*/

func testCreateWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	arg := &warehouse.CreateWarehouseArg{
		Body: &warehouse.Warehouse{
			Name:        proto.String("Buck1"),
			Description: proto.String("Docker Sawrm."),
			Metadata:    nil,
		},
	}
	glog.Infof("Warehouse Arg %v", arg)
	response, err := client.CreateWarehouse(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Warehouse Response %s", response)
	// taskUUID = response.GetContent().GetTaskReferenceData().GetValue().GetExtId()
	// client.GetTaskReferenceData()

}

func testUpdateWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	arg := &warehouse.UpdateWarehouseMetadataArg{
		ExtId: proto.String("7db350a6-3e18-4723-5a0f-a5c076bc17b2"),
		Body: &warehouse.Warehouse{
			Name:        proto.String("Mast & Harbor"),
			Description: proto.String("Flipkart"),
			Metadata:    nil,
		},
	}

	glog.Infof("Warehouse Arg %v", arg)
	response, err := client.UpdateWarehouseMetadata(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Warehouse Response ", response)

}

func testSyncWarehouse(ctx context.Context, client warehouse.WarehouseServiceClient) {
	arg := &warehouse.SyncWarehouseMetadataArg{
		ExtId: proto.String("ec1125d9-0ec5-4e03-4535-fc67a1b7cda9"),
	}

	glog.Infof("Warehouse Arg %v", arg)
	response, err := client.SyncWarehouseMetadata(ctx, arg)
	if err != nil {
		fmt.Println("Marina request error:", err)
		glog.Errorf("Marina request error: %s", err)
		return
	}

	fmt.Println("***********Response received from server*******")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Warehouse Response ", response)

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func main() {
	// TODO: Consider moving these to init().
	var grpcServerPort uint
	grpcServerPort = 9200 // 32391 //9200 //30188 //9200
	// glog.Infof("gRPC server port to connect: %v", grpcServerPort)
	// fmt.Println("gRPC server port to connect: ", grpcServerPort)
	conn, whclient := warehouseClient(grpcServerPort)
	conn, wiclient := warehouseItemClient(grpcServerPort)
	if wiclient == nil || whclient == nil {
		glog.Errorf("Error occurred in creating Clients %s %s", whclient, wiclient)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	/*           WarehouseGet                     */
	// fmt.Println("Making Marina server request WarehouseGet.")
	// testGetWarehouse(ctx, whclient)

	/*           WarehouseCreate                     */
	// fmt.Println("Creating Marina Warehouse.")
	// testCreateWarehouse(ctx, whclient)

	/*           WarehouseList                     */
	// fmt.Println("Making Marina server request WarehouseGet.")
	// testListWarehouse(ctx, whclient)

	// /*           WarehouseDelete                     */
	// fmt.Println("Deleting Marina Warehouse.")
	// testDeleteWarehouse(ctx, whclient)

	/*           WarehouseUpdate                     */
	// fmt.Println("Updating Marina Warehouse.")
	// testUpdateWarehouse(ctx, whclient)
	// Fetching after update
	/*           WarehouseGet                     */
	// time.Sleep(3 * time.Second)
	// fmt.Println("Making Marina server request WarehouseGet.")
	// testGetWarehouse(ctx, whclient)
	// /*           WarehouseSync                    */
	fmt.Println("Syncing Marina Warehouse.")
	testSyncWarehouse(ctx, whclient)

	//     WarehouseItems Operations                  */
	/*           AddItemToWarehouse                     */
	// fmt.Println("Making Marina server request AddItemToWarehouse.")
	// testCreateWarehouseItem(ctx, wiclient)

	/*           GetWarehouseItemById in a Warehouse                     */
	// fmt.Println("Making Marina server request GetWarehouseItemId.")
	// testGetWarehouseItemById(ctx, wiclient)

	/*           ListWarehouseItems in a Warehouse                     */
	// fmt.Println("Making Marina server request ListWarehouseItems in Warehouse.")
	// testListWarehouseItem(ctx, wiclient)

	/*           DeleteWarehouseItemById in a Warehouse                     */
	// fmt.Println("Making Marina server request DeleteWarehouseItem .")
	// testDeleteWarehouseItemById(ctx, wiclient)

	/*           UpdateWarehouseItemById in a Warehouse                     */
	// fmt.Println("Making Marina server request UpdateWarehouseItem .")
	// testUpdateWarehouseItemById(ctx, wiclient)

	// Delete all warehouse items and warehouses from IDF.
	// testDeleteAllWarehouseItems(ctx, wclient)
	// deleteAllWarehouses(ctx, whclient)

}
