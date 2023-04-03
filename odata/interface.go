/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 */

package odata

import (
	"context"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/content-management-marina/authz"
	"github.com/nutanix-core/content-management-marina/db"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

type OdataInterface interface {
	GetLimitAndOffset(ctx context.Context, idfParams *marinaIfc.IdfParameters) (*int64, *int64, error)
	QueryWithAuthorizedUuids(ctx context.Context, idfIfc db.IdfClientInterface, authorized *authz.Authorized,
		argUuids [][]byte, idfParams *marinaIfc.IdfParameters, queryName string, entityType db.EntityType,
		tableCols ...interface{}) ([]*insights_interface.EntityWithMetric, int64, error)
}
