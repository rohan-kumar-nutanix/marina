/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package security_policy

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/db"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type ISecurityPolicyDB interface {
	CreateSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, securityPolicyUuid *uuid4.Uuid, securityPolicyPB *config.SecurityPolicy) error
	GetSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, securityPolicyUuid *uuid4.Uuid) (
		*config.SecurityPolicy, error)
	DeleteSecurityPolicy(ctx context.Context, idfIfc db.IdfClientInterface,
		cpdbIfc cpdb.CPDBClientInterface, securityPolicyUuid string) error
	UpdateSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, securityPolicyUuid *uuid4.Uuid, securityPolicyPB *config.SecurityPolicy) error
	ListSecurityPolices(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) ([]*config.SecurityPolicy, error)
}
