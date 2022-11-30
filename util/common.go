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
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

/*
itemType - Catalog Item type in string.
Returns CatalogItemInfo Enum Type.
*/
func GetCatalogItemTypeEnum(itemType string) *marinaIfc.CatalogItemInfo_CatalogItemType {
	switch itemType {
	case "kImage":
		return marinaIfc.CatalogItemInfo_kImage.Enum()
	case "kAcropolisVmSnapshot":
		return marinaIfc.CatalogItemInfo_kAcropolisVmSnapshot.Enum()
	case "kVmSnapshot":
		return marinaIfc.CatalogItemInfo_kVmSnapshot.Enum()
	case "kFile":
		return marinaIfc.CatalogItemInfo_kFile.Enum()
	case "kLCM":
		return marinaIfc.CatalogItemInfo_kLCM.Enum()
	case "kOVA":
		return marinaIfc.CatalogItemInfo_kOVA.Enum()
	case "kVmTemplate":
		return marinaIfc.CatalogItemInfo_kVmTemplate.Enum()
	}
	return nil
}
