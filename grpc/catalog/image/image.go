/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
 * Wrapper around the Image DB entry. Includes libraries that will
 * interact with IDF and query for Images.
*/

package image

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nutanix-core/acs-aos-go/acropolis"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	utils "github.com/nutanix-core/content-management-marina/util"
)

var imageInterface ImageInterface = nil
var once sync.Once

type ImageImpl struct {
}

func NewImageInterface() ImageInterface {
	once.Do(func() {
		imageInterface = new(ImageImpl)
	})
	return imageInterface
}

func (*ImageImpl) CreateImage(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, protoIfc utils.ProtoUtilInterface,
	imageUuid *uuid4.Uuid, image *acropolis.ImageInfo) error {

	attrParams := make(cpdb.AttributeParams)
	attrParams["uuid"] = imageUuid.String()
	attrParams["name"] = image.GetName()
	attrParams["annotation"] = image.GetAnnotation()
	attrParams["image_type"] = image.GetImageType().String()
	attrParams["image_state"] = image.GetImageState().String()
	// TODO: Add "hypervisor" field after updating acropolis protos

	if image.Hidden != nil {
		attrParams["hidden"] = image.GetHidden()
	}

	if image.Architecture != nil {
		attrParams["architecture"] = image.GetArchitecture().String()
	}

	if image.SourceInsecureUrl != nil {
		attrParams["source_insecure_url"] = image.GetSourceInsecureUrl()
	}

	if image.ContainerId != nil {
		attrParams["container_id"] = image.GetContainerId()
	}

	if image.Deleted != nil {
		attrParams["deleted"] = image.GetDeleted()
	}

	if image.ChecksumType != nil {
		attrParams["checksum_type"] = image.GetChecksumType().String()
	}

	if image.ChecksumBytes != nil {
		attrParams["checksum_bytes"] = image.GetChecksumBytes()
	}

	if image.LogicalTimestamp != nil {
		attrParams["logical_timestamp"] = image.GetLogicalTimestamp()
	}

	if image.VmdiskUuid != nil {
		attrParams["vmdisk_uuid"] = uuid4.ToUuid4(image.GetVmdiskUuid()).String()
	}

	if image.CreateTimeUsecs != nil {
		attrParams["create_time_usecs"] = image.GetCreateTimeUsecs()
	}

	if image.UpdateTimeUsecs != nil {
		attrParams["update_time_usecs"] = image.GetUpdateTimeUsecs()
	}

	if image.VmdiskSize != nil {
		attrParams["vmdisk_size"] = image.GetVmdiskSize()
	}

	if image.VmdiskNfsPath != nil {
		attrParams["vmdisk_nfs_path"] = image.GetVmdiskNfsPath()
	}

	if image.InCatalog != nil {
		attrParams["in_catalog"] = image.GetInCatalog()
	}

	if image.FileUuid != nil {
		attrParams["file_uuid"] = uuid4.ToUuid4(image.GetFileUuid()).String()
	}

	if image.OwnerClusterUuid != nil {
		attrParams["owner_cluster_uuid"] = uuid4.ToUuid4(image.GetOwnerClusterUuid()).String()
	}

	if image.NewImageUuid != nil {
		attrParams["new_image_uuid"] = uuid4.ToUuid4(image.GetNewImageUuid()).String()
	}

	if image.OldImageUuid != nil {
		attrParams["old_image_uuid"] = uuid4.ToUuid4(image.GetOldImageUuid()).String()
	}

	if image.InitialClusterLocationList != nil {
		var initialClusters []string
		for _, cluster := range image.GetInitialClusterLocationList() {
			initialClusters = append(initialClusters, uuid4.ToUuid4(cluster).String())
		}
		attrParams["initial_cluster_location_list"] = initialClusters
	}

	marshal, err := protoIfc.Marshal(image)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the image proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the image proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	err = writer.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while closing zlib stream: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = buffer.Bytes()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while generating IDF attributes: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	_, err = cpdbIfc.UpdateEntity(db.Image.ToString(), imageUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for image %s: %v", imageUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}
