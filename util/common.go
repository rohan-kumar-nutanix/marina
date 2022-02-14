/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * This file implements utility functions to do localhost operations when the
 * service is running inside a docker container.
 */

package utils

import (
	marina_pb "github.com/nutanix-core/content-management-marina/protos/marina"
)

/*
	itemType - Catalog Item type in string.
	Returns CatalogItemInfo Enum Type.
*/
func GetCatalogItemTypeEnum(itemType string) *marina_pb.CatalogItemInfo_CatalogItemType {
	switch itemType {
	case "kImage":
		return marina_pb.CatalogItemInfo_kImage.Enum()
	case "kAcropolisVmSnapshot":
		return marina_pb.CatalogItemInfo_kAcropolisVmSnapshot.Enum()
	case "kVmSnapshot":
		return marina_pb.CatalogItemInfo_kVmSnapshot.Enum()
	case "kFile":
		return marina_pb.CatalogItemInfo_kFile.Enum()
	case "kLCM":
		return marina_pb.CatalogItemInfo_kLCM.Enum()
	case "kOVA":
		return marina_pb.CatalogItemInfo_kOVA.Enum()
	case "kVmTemplate":
		return marina_pb.CatalogItemInfo_kVmTemplate.Enum()
	}
	return nil
}
