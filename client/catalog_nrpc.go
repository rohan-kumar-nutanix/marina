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

	catalogIfc "github.com/nutanix-core/acs-aos-go/catalog"
	catalogClient "github.com/nutanix-core/acs-aos-go/catalog/client"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
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
			//uuid4.ToUuid4(t.TemplateInfo.TemplateUuid).UuidToString()
			fmt.Println(*t.TemplateInfo.Name, len(t.VersionInfoList))
			fmt.Println("***Template Object ***: ", t)
			time.Sleep(1 * time.Second)
		}
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
			// fmt.Println("'%v': %v\n", i, item)
			fmt.Println(uuid4.ToUuid4(item.GlobalCatalogItemUuid).UuidToString(), *item.Name, "\t", *item.ItemType, "\t", *item.Version)
		}
	}

}

func catalogItemsGet() {
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
			// fmt.Println("'%v': %v\n", i, item)
			fmt.Println(uuid4.ToUuid4(item.GlobalCatalogItemUuid).UuidToString(), *item.Name, "\t", *item.ItemType, "\t", *item.Version)
		}
	}

}

func filesGet() {
	fmt.Println("\n\n\n-----------------------Fetching Files-------------------")
	fmt.Println("-----------------------------------------------------------------------------------------------------\n\n")
	filesGetArg := &catalogIfc.FileGetArg{}
	filesGetRet, err := catalogRpcClient.FileGet(filesGetArg)
	if err != nil {
		fmt.Println("Error in fetching CatalogItemDetails :", err)
	} else {
		fmt.Println("F_UUID                                   Name   ")
		fmt.Println("-----------------------------------------------------------------------------------------------------")
		for _, f := range filesGetRet.FileInfoList {
			//uuid4.ToUuid4(t.TemplateInfo.TemplateUuid).UuidToString()
			fmt.Println("***File Object Details ***: ", f)
		}
	}

}

func main() {

	catalogRpcClient = catalogClient.NewCatalogService("localhost", 9202)
	catalogSvcUtil = clientUtil.NewCatalogService("localhost", 9202)

	// Fetch CatalogItems.
	fmt.Println("\n********************************************************")
	fmt.Println("Fetching CatalogItems....")
	catalogItemsGet()
	fmt.Println("Fetching CatalogItems using SendMsg method....")
	catalogItemsGetWithSendMsg()

	fmt.Println("\n********************************************************")
	time.Sleep(3 * time.Second)
	fmt.Println("Fetching VMTemplates....")
	testVmTemplates()

	time.Sleep(3 * time.Second)
	fmt.Println("\n********************************************************")
	fmt.Println("Fetching FileItems....")
	filesGet()

}
