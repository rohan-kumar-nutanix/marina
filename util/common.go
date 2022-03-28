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
	marinapb "github.com/nutanix-core/content-management-marina/protos/marina"
)

/*
	itemType - Catalog Item type in string.
	Returns CatalogItemInfo Enum Type.
*/
func GetCatalogItemTypeEnum(itemType string) *marinapb.CatalogItemInfo_CatalogItemType {
	switch itemType {
	case "kImage":
		return marinapb.CatalogItemInfo_kImage.Enum()
	case "kAcropolisVmSnapshot":
		return marinapb.CatalogItemInfo_kAcropolisVmSnapshot.Enum()
	case "kVmSnapshot":
		return marinapb.CatalogItemInfo_kVmSnapshot.Enum()
	case "kFile":
		return marinapb.CatalogItemInfo_kFile.Enum()
	case "kLCM":
		return marinapb.CatalogItemInfo_kLCM.Enum()
	case "kOVA":
		return marinapb.CatalogItemInfo_kOVA.Enum()
	case "kVmTemplate":
		return marinapb.CatalogItemInfo_kVmTemplate.Enum()
	}
	return nil
}
