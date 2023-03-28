/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * File Repo interface
 *
 */

package file_repo

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type FileRepoInterface interface {
	GetFile(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, fileUuid *uuid4.Uuid) (*marinaIfc.FileInfo, error)
	GetFiles(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface, fileUuids []*uuid4.Uuid) (
		[]*marinaIfc.FileInfo, error)
	CreateFile(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, protoIfc utils.ProtoUtilInterface,
		fileUuid *uuid4.Uuid, peUuid *uuid4.Uuid, peFileUuid *uuid4.Uuid) error
}
