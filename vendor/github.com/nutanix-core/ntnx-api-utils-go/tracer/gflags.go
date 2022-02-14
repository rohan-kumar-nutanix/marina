//
// Copyright (c) 2021 Nutanix Inc. All rights reserved.
//
// Author: deepanshu.singhal@nutanix.com
//

package tracer

import (
	"flag"
)

var observabilityConfigFilePath = flag.String(
	"observability_config_file_path",
	"/home/nutanix/config/observability/config.json",
	"Path for file with observability configuration")
