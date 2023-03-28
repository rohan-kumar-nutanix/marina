/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Gflags defined for Catalog Item module.
 */

package catalog_item

import (
	"flag"
	"fmt"

	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

var CatalogIdfQueryChunkSize = flag.Int("catalog_idf_query_chunk_size", 30,
	"Maximum number of catalog items that can be queried in one IDF query")

var RegistrationTaskRetryDelayMs = flag.Int("registration_task_retry_delay_ms", 5*1000,
	"Delay in milliseconds before next attempt.")

var CatalogCreatSpecItemTypeList = flag.String("catalog_create_spec_item_type_list",
	fmt.Sprintf("%v,%v,%v", marinaIfc.CatalogItemInfo_kImage, marinaIfc.CatalogItemInfo_kFile,
		marinaIfc.CatalogItemInfo_kLCM),
	"Comma separated list of entity types for which PC owned spec needs to be created.")
