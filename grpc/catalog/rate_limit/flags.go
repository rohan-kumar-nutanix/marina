/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Gflags defined for Rate Limit module.
 */

package rate_limit

import "flag"

var catalogRateLimitCategoriesCacheSize = flag.Int("catalog_rate_limit_categories_cache_size", 32,
	"Maximum size of rate limits categories cache")

var catalogRateLimitCategoriesCacheTimeoutSecs = flag.Int("catalog_rate_limit_categories_cache_timeout_secs", 5,
	"How long are rate limits categories cache items considered valid")
