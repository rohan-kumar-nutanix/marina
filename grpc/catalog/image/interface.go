/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * File Repo interface
 *
 */

package image

import (
	"context"

	"github.com/nutanix-core/acs-aos-go/acropolis"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type ImageInterface interface {
	CreateImage(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, protoIfc utils.ProtoUtilInterface,
		imageUuid *uuid4.Uuid, image *acropolis.ImageInfo) error
}
