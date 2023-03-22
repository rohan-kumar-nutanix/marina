/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
 * Wrapper around the File Repo DB entry. Includes libraries that will
 * interact with IDF and query for Files.
*/

package file_repo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

var fileRepoInterface FileRepoInterface = nil
var once sync.Once

type FileRepoImpl struct {
}

func NewFileRepoImpl() FileRepoInterface {
	once.Do(func() {
		fileRepoInterface = new(FileRepoImpl)
	})
	return fileRepoInterface
}

func (*FileRepoImpl) GetFile(cpdbIfc cpdb.CPDBClientInterface, fileUuid *uuid4.Uuid) (*marinaIfc.FileInfo, error) {
	entity, err := cpdbIfc.GetEntity(db.File.ToString(), fileUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("File %s not found", fileUuid.String())
		return nil, marinaError.ErrNotFound

	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting File %s", fileUuid.String())
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))

	}

	file := &marinaIfc.FileInfo{}
	err = entity.DeserializeEntity(file)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deserialize file IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return file, nil
}

func (*FileRepoImpl) GetFiles(cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface,
	fileUuids []*uuid4.Uuid) ([]*marinaIfc.FileInfo, error) {

	var guids []*insights_interface.EntityGuid
	for _, uuid := range fileUuids {
		guids = append(guids, &insights_interface.EntityGuid{
			EntityTypeName: proto.String(db.File.ToString()),
			EntityId:       proto.String(uuid.String()),
		})
	}

	entities, err := cpdbIfc.GetEntities(guids, false)
	if err == insights_interface.ErrNotFound {
		log.Errorf("Files %s do not exist", uuidIfc.StringFromUuidPointers(fileUuids))
		return []*marinaIfc.FileInfo{}, nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Files %s", uuidIfc.StringFromUuidPointers(fileUuids))
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))

	}

	var files []*marinaIfc.FileInfo
	for _, entity := range entities {
		file := &marinaIfc.FileInfo{}
		err = entity.DeserializeEntity(file)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to deserialize file IDF entry: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		files = append(files, file)
	}

	return files, nil
}
