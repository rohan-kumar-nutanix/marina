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
)

var CatalogIdfQueryChunkSize = flag.Int("catalog_idf_query_chunk_size", 30,
	"Maximum number of catalog items that can be queried in one IDF query")
