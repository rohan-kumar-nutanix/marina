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
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type FileRepoInterface interface {
	GetFile(cpdbIfc cpdb.CPDBClientInterface, fileUuid *uuid4.Uuid) (*marinaIfc.FileInfo, error)
	GetFiles(cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface, fileUuids []*uuid4.Uuid) (
		[]*marinaIfc.FileInfo, error)
}
