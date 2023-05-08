/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package scanner_config

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/db"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type IScannerConfigDB interface {
	CreateScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, scannerConfigUuid *uuid4.Uuid, scannerConfigPB *config.ScannerConfig) error
	GetScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, scannerConfigUuid *uuid4.Uuid) (
		*config.ScannerConfig, error)
	DeleteScannerConfig(ctx context.Context, idfIfc db.IdfClientInterface,
		cpdbIfc cpdb.CPDBClientInterface, scannerConfigUuid string) error
	UpdateScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, scannerConfigUuid *uuid4.Uuid, scannerConfigPB *config.ScannerConfig) error
	ListScannerConfigurations(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) ([]*config.ScannerConfig, error)
}
