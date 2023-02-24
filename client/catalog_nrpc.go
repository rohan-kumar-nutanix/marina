/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * This file is make tests Catalog RPC requests.
 *
 */

package main

import (
	"fmt"
	"time"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	catalogClient "github.com/nutanix-core/acs-aos-go/catalog/client"
	catalogIfc "github.com/nutanix-core/acs-aos-go/catalog_pc"

	clientUtil "github.com/nutanix-core/content-management-marina/util/catalog/client"
)

var catalogSvcUtil *clientUtil.Catalog
var catalogRpcClient *catalogClient.Catalog

func testVmTemplates() {
	fmt.Println("\n\n-------------------Fetching VMTemplates-------------------")
	fmt.Println("-----------------------------------------------------------------------------------------------------")

	vmTemplatesGetArg := &catalogIfc.TemplatesGetArg{}
	vmTemplatesGetRet := &catalogIfc.TemplatesGetRet{}

	err := catalogSvcUtil.SendMsg("VmTemplatesGet", vmTemplatesGetArg, vmTemplatesGetRet)
	if err != nil {
		fmt.Println("Error in fetching VMTemplate Details :", err)
	} else {
		fmt.Printf("vmTemplatesGetRet: %v \n", vmTemplatesGetRet)
		fmt.Println("VMTemplates got retrived total: ", len(vmTemplatesGetRet.TemplateDetailsList))
		fmt.Println("T_UUID                                   Name             ActiveVersion              No_Versions")
		fmt.Println("-----------------------------------------------------------------------------------------------------")
		for _, t := range vmTemplatesGetRet.TemplateDetailsList {
			// uuid4.ToUuid4(t.TemplateInfo.TemplateUuid).UuidToString()
			// fmt.Println(*t.TemplateInfo.Name, len(t.VersionInfoList))
			fmt.Println("***Template Object ***: ", t)
			time.Sleep(1 * time.Second)
		}
	}

}

func createTemplateShell(name, desc *string) {
	arg := &catalogIfc.TemplateCreateArg{}
	arg.TemplateName = name
	arg.TemplateDescription = desc
	ret := &catalogIfc.TemplateCreateRet{}
	err := catalogSvcUtil.SendMsg("VmTemplateCreate", arg, ret)
	if err != nil {
		fmt.Println("Error in Creating VMTemplate Shell :", err)
	} else {
		fmt.Println("Successfully  VMTemplate Shell got created Task uuid :", uuid4.ToUuid4(ret.TaskUuid).String())
	}
}
func createVMTemplate(name, desc, vmuuid *string) {
	arg := &catalogIfc.TemplateAndVersionCreateArg{}
	arg.TemplateName = name
	arg.TemplateDescription = desc
	vmUuid, err := uuid4.StringToUuid4(*vmuuid)
	if err == nil {
		arg.VmUuid = vmUuid.RawBytes()
	} else {
		fmt.Println("Invalid VMUUID passed to function. Returning.. err ", err)
		return
	}

	ret := &catalogIfc.TemplateAndVersionCreateRet{}
	err = catalogSvcUtil.SendMsg("VmTemplateAndVersionCreate", arg, ret)
	if err != nil {
		fmt.Println("Error in Creating VMTemplate and version :", err)
	} else {
		fmt.Println("Successfully  VMTemplate and version got created Task uuid :", uuid4.ToUuid4(ret.TaskUuid).String())
	}
}

func catalogItemsGetWithSendMsg() {
	cItemGetArg := &catalogIfc.CatalogItemGetArg{}
	cItemGetRet := &catalogIfc.CatalogItemGetRet{}
	err := catalogSvcUtil.SendMsg("CatalogItemGet", cItemGetArg, cItemGetRet)
	if err != nil {
		fmt.Println("Error in fetching CatalogItemDetails :", err)
	} else {
		fmt.Println("Catalog Items got retrived total: ", len(cItemGetRet.CatalogItemList))

		fmt.Printf("********Legacy Catalog Items********\n-------------------------\n")
		fmt.Println("G_UUID                                   Name           		Type              Version")
		fmt.Println("-----------------------------------------------------------------------------------------------------")
		for _, item := range cItemGetRet.CatalogItemList {
			fmt.Printf("%v\n", item)
			// fmt.Println(uuid4.ToUuid4(item.GlobalCatalogItemUuid).UuidToString(), *item.Name, "\t", *item.ItemType, "\t", *item.Version)
		}
	}

}

/*func catalogItemsGet() {
	cItemGetArg := &catalogIfc.CatalogItemGetArg{}
	cItemGetRet, err := catalogRpcClient.CatalogItemGet(cItemGetArg)
	if err != nil {
		fmt.Println("Error in fetching CatalogItemDetails :", err)
	} else {
		fmt.Println("Catalog Items got retrived total: ", len(cItemGetRet.CatalogItemList))

		fmt.Printf("********Legacy Catalog Items********\n-------------------------\n")
		fmt.Println("G_UUID                                   Name            Type              Version")
		fmt.Println("-----------------------------------------------------------------------------------------------------")
		for _, item := range cItemGetRet.CatalogItemList {
			fmt.Println(item)
			// fmt.Println(uuid4.ToUuid4(item.GlobalCatalogItemUuid).UuidToString(), *item.Name, "\t", *item.ItemType, "\t", *item.Version)
		}
	}

}*/

/*func catalogItemsCreate(name, desc , imageUrl *string, ) {
	arg := &catalogIfc.CatalogItemCreateArg{}
	spec := &catalogIfc.CatalogItemCreateSpec{}
	spec.Name = name
	spec.Annotation = desc
	sourceGroup := &catalogIfc.SourceGroupSpec{}
	sourceGroup.
	append(a, )
	spec.
	catalogItem := &catalogIfc.CatalogItemId{}
	catalogItem.GlobalCatalogItemUuid = uuid4.EmptyUuid().RawBytes()
	cItemDelArg.CatalogItemId = catalogItem
	cItemDelRet, err := catalogRpcClient.CatalogItemDelete(cItemDelArg)
	if err != nil {
		fmt.Println("Error in Deleting CatalogItemD :", err)
	} else {
		fmt.Println("*** Successfully Deleted CatalogItemD :", uuid4.ToUuid4(cItemDelRet.GetTaskUuid()).String())
	}
}*/

/*func catalogItemsDelete() {
	cItemDelArg := &catalogIfc.CatalogItemDeleteArg{}
	arg := &acropolisIfc.ImageCreateArg{}
	spec1 := &acropolisIfc.ImageCreateSpec{}
	spec1.

	catalogItem := &catalogIfc.CatalogItemId{}
	catalogItem.GlobalCatalogItemUuid = uuid4.EmptyUuid().RawBytes()
	cItemDelArg.CatalogItemId = catalogItem
	cItemDelRet, err := catalogRpcClient.CatalogItemDelete(cItemDelArg)

	if err != nil {
		fmt.Println("Error in Deleting CatalogItemD :", err)
	} else {
		fmt.Println("*** Successfully Deleted CatalogItemD :", uuid4.ToUuid4(cItemDelRet.GetTaskUuid()).String())
	}
}*/

/*func filesGet() {
	fmt.Println("\n\n\n-----------------------Fetching Files-------------------")
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	filesGetArg := &catalogIfc.FileGetArg{}
	filesGetRet, err := catalogRpcClient.FileGet(filesGetArg)
	if err != nil {
		fmt.Println("Error in fetching CatalogItemDetails :", err)
	} else {
		fmt.Println("F_UUID                                   Name   ")
		fmt.Println("-----------------------------------------------------------------------------------------------------")
		for _, f := range filesGetRet.FileInfoList {
			// uuid4.ToUuid4(t.TemplateInfo.TemplateUuid).UuidToString()
			fmt.Println("***File Object Details ***: ", f)
		}
	}

}*/

func main() {

	catalogSvcUtil = clientUtil.NewCatalogService("localhost", 9202)
	// Fetch CatalogItems.
	fmt.Println("\n********************************************************")
	fmt.Println("Fetching CatalogItems....")
	// catalogItemsGet()
	// fmt.Println("Fetching CatalogItems using SendMsg method....")
	// catalogItemsGetWithSendMsg()

	// Deleting CatalogItems.
	// fmt.Println("Deleting CatalogItems....")
	// catalogItemsDelete()
	var tname string = "temp1"
	var desc string = "my Temp Desc"
	// var vmUuid string = "1bbd7143-e127-4d17-bc52-b429d6e99ab4"
	fmt.Println("Cteating VMTemplate Shell......")
	createTemplateShell(&tname, &desc)
	/*fmt.Println("Cteating VMTemplate and Version vm uuid= ", vmUuid)
	createVMTemplate(&tname, &desc, &vmUuid)*/
	fmt.Println("\n********************************************************")
	time.Sleep(3 * time.Second)
	fmt.Println("Fetching VMTemplates....")
	testVmTemplates()

	/*	time.Sleep(3 * time.Second)
		fmt.Println("\n********************************************************")
		fmt.Println("Fetching FileItems....")*/
	// filesGet()

}
