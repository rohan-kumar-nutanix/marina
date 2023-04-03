/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 */

package rate_limit

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	cache "github.com/go-pkgz/expirable-cache/v2"
	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	categoryUtil "github.com/nutanix-core/acs-aos-go/aplos/categories/category_util"
	filterUtil "github.com/nutanix-core/acs-aos-go/aplos/categories/filter_util"
	aplosProto "github.com/nutanix-core/acs-aos-go/aplos/categories/proto"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/authz"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/metadata"
	"github.com/nutanix-core/content-management-marina/odata"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/content-management-marina/zeus"
)

const (
	redacted                = "-"
	rateLimitIamObjectType  = "image_bw_throttling_policy"
	viewRateLimitPermission = "CatalogService:View_Image_BW_Throttling_Policy"

	clusterKind         = "cluster"
	rateLimitFilterKind = "image_rate_limit"

	kindCol           = "kind"
	kindIdCol         = "kind_id"
	categoryIdListCol = "category_id_list"
)

var rateLimitFilterStringToEnum = map[string]marinaIfc.CatalogRateLimitFilter_CategoryMatchType{
	"CATEGORIES_MATCH_ALL": marinaIfc.CatalogRateLimitFilter_kAll,
	"CATEGORIES_MATCH_ANY": marinaIfc.CatalogRateLimitFilter_kAny,
}

var catalogRateLimitAttribute = []interface{}{
	insights_interface.COMPRESSED_PROTOBUF_ATTR,
}

var ecapAttributes = []interface{}{
	categoryIdListCol,
}

var (
	// Caches the return values of GetCategoryDetailsFromCategoryUuids, which resolves the categories names for
	// category UUIDs
	categoryNamesCache     cache.Cache[string, map[string][]string]
	categoryNamesCacheOnce sync.Once

	// Caches the return values of GetEntityUuidsByKindForCategoryUuids, which resolves matched entities for category
	// UUIDs and match type
	filterEntitiesCache     cache.Cache[string, []*uuid4.Uuid]
	filterEntitiesCacheOnce sync.Once
)

// CatalogRateLimitGet implements the CatalogRateLimitGet RPC.
func CatalogRateLimitGet(ctx context.Context, arg *marinaIfc.CatalogRateLimitGetArg, idfIfc db.IdfClientInterface,
	cpdbIfc cpdb.CPDBClientInterface, authzIfc authz.AuthzInterface, iamIfc authz_cache.IamClientIfc,
	metadataIfc metadata.EntityMetadataInterface, odataIfc odata.OdataInterface,
	categoryIfc categoryUtil.CategoryUtilInterface, filterIfc filterUtil.FilterUtilInterface, zeusConfig zeus.ConfigCache) (
	*marinaIfc.CatalogRateLimitGetRet, error) {

	log.V(2).Info("CatalogRateLimit RPC started.")
	defer log.V(2).Info("CatalogRateLimit RPC finished.")

	authorized, err := authzIfc.GetAuthorizedUuids(ctx, iamIfc, cpdbIfc, rateLimitIamObjectType,
		db.RateLimit.ToString(), viewRateLimitPermission)
	if err != nil {
		return nil, err
	}

	entities, count, err := odataIfc.QueryWithAuthorizedUuids(ctx, idfIfc, authorized, arg.GetRateLimitUuidList(),
		arg.GetIdfParams(), "catalog_rate_limit_get:get_odata", db.RateLimit, catalogRateLimitAttribute...)
	if err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return &marinaIfc.CatalogRateLimitGetRet{TotalEntityCount: proto.Int64(count)}, nil
	}

	var rateLimits []*marinaIfc.CatalogRateLimitInfo
	for _, entityWithMetric := range entities {
		rateLimit := &marinaIfc.CatalogRateLimitInfo{}
		err = entityWithMetric.DeserializeEntity(rateLimit)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to deserialize rate limit IDF entry: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		rateLimits = append(rateLimits, rateLimit)
	}

	if arg.GetIncludeMetadata() {
		var rateLimitUuids []string
		for _, rateLimit := range rateLimits {
			rateLimitUuids = append(rateLimitUuids, uuid4.ToUuid4(rateLimit.Uuid).String())
		}

		metadataByUuid, err := metadataIfc.GetEntityMetadata(ctx, cpdbIfc, rateLimitUuids, db.RateLimit.ToString(),
			"catalog_rate_limit_get:get_metadata")
		if err != nil {
			return nil, err
		}

		for _, rateLimit := range rateLimits {
			rateLimit.Metadata = metadataByUuid[*uuid4.ToUuid4(rateLimit.Uuid)]
		}
	}

	if arg.GetIncludeStatus() {
		for _, rateLimit := range rateLimits {
			err = fillStatus(ctx, cpdbIfc, categoryIfc, filterIfc, zeusConfig, rateLimit, authorized)
			if err != nil {
				return nil, err
			}
		}
	}

	ret := &marinaIfc.CatalogRateLimitGetRet{
		RateLimitInfoList: rateLimits,
		TotalEntityCount:  proto.Int64(count),
	}

	// Add E-Tag to the response if queried for only one Rate Limit, and the include_etag flag was True
	if arg.GetIncludeEtag() && len(ret.GetRateLimitInfoList()) == 1 {
		etag := strconv.FormatInt(ret.GetRateLimitInfoList()[0].GetLogicalTimestamp(), 10)
		ret.Etag = proto.String(etag)
	}

	return ret, nil
}

func getCategoryNamesCache() cache.Cache[string, map[string][]string] {
	categoryNamesCacheOnce.Do(func() {
		categoryNamesCache = cache.NewCache[string, map[string][]string]().WithLRU().
			WithMaxKeys(*catalogRateLimitCategoriesCacheSize).
			WithTTL(time.Second * time.Duration(*catalogRateLimitCategoriesCacheTimeoutSecs))
	})
	return categoryNamesCache
}

func getFilterEntitiesCache() cache.Cache[string, []*uuid4.Uuid] {
	filterEntitiesCacheOnce.Do(func() {
		filterEntitiesCache = cache.NewCache[string, []*uuid4.Uuid]().WithLRU().
			WithMaxKeys(*catalogRateLimitCategoriesCacheSize).
			WithTTL(time.Second * time.Duration(*catalogRateLimitCategoriesCacheTimeoutSecs))
	})
	return filterEntitiesCache
}

func fillStatus(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, categoryIfc categoryUtil.CategoryUtilInterface,
	filterIfc filterUtil.FilterUtilInterface, zeusConfig zeus.ConfigCache, rateLimit *marinaIfc.CatalogRateLimitInfo,
	authorized *authz.Authorized) error {

	filterUuid := uuid4.ToUuid4(rateLimit.GetClusterFilterUuid())
	filter, err := filterIfc.GetFilter(ctx, filterUuid)
	if err != nil {
		errMsg := fmt.Sprintf("Filter %s not found", filterUuid.String())
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	status := &marinaIfc.CatalogRateLimitInfo_CatalogRateLimitStatusInfo{}
	for _, filterExpr := range filter.GetFilterExpressions() {
		if filterExpr.GetLhsEntityType() == "category" &&
			filterExpr.GetExpressionType() == aplosProto.FilterExpression_kScope {

			// Fill the filter status
			filterMatchType := filterUtil.FilterMatchTypeEnumToString[filterExpr.GetOperator()]
			rateLimitFilter := rateLimitFilterStringToEnum[filterMatchType]
			status.ClusterFilter = &marinaIfc.CatalogRateLimitFilter{
				MatchType: &rateLimitFilter,
			}

			var categoryCacheKey string
			var categoryUuids []*uuid4.Uuid
			categoryUuidStrs := strings.Split(filterExpr.GetEntityUuids(), ":")
			sort.Strings(categoryUuidStrs)
			for _, categoryUuidStr := range categoryUuidStrs {
				categoryUuid, err := uuid4.StringToUuid4(categoryUuidStr)
				if err != nil {
					errMsg := fmt.Sprintf("Error encountered while creating category UUID: %v", err)
					return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
				}
				categoryUuids = append(categoryUuids, categoryUuid)
				categoryCacheKey = categoryCacheKey + categoryUuidStr
			}

			var categories map[string][]string
			categories, ok := getCategoryNamesCache().Get(categoryCacheKey)
			if !ok {
				categories, err = categoryIfc.GetCategoryDetailsFromCategoryUuids(ctx, categoryUuids)
				if err != nil {
					errMsg := fmt.Sprintf("Error encountered while getting category details: %v", err)
					return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
				}
				getCategoryNamesCache().Set(categoryCacheKey, categories, 0 /* Cache wide TTl will be used */)
			}

			for categoryKey, categoryValues := range categories {
				status.GetClusterFilter().Categories = append(status.GetClusterFilter().Categories,
					&marinaIfc.CatalogRateLimitFilter_CategoryKeyValue{
						CategoryKey:    proto.String(categoryKey),
						CategoryValues: categoryValues,
					},
				)
			}

			// Fill the cluster status
			var clusterUuids []*uuid4.Uuid
			filterCacheKey := categoryCacheKey + filterMatchType
			clusterUuids, ok = getFilterEntitiesCache().Get(filterCacheKey)
			if !ok {
				clusterUuids, err = categoryIfc.GetEntityUuidsByKindForCategoryUuids(ctx, []string{clusterKind},
					categoryUuids, filterMatchType)
				if err != nil && err != insights_interface.ErrNotFound {
					errMsg := fmt.Sprintf("Error encountered while getting cluster UUIDs: %v", err)
					return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
				}
				getFilterEntitiesCache().Set(filterCacheKey, clusterUuids, 0 /* Cache wide TTl will be used */)
			}

			for _, clusterUuid := range clusterUuids {
				isSupported := zeusConfig.IsRateLimitSupported(clusterUuid)
				if isSupported == nil {
					errMsg := "invalid zeus config. IsRateLimitSupported property not found"
					return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
				}

				if !*isSupported {
					status.UnsupportedClusterUuidList = append(status.UnsupportedClusterUuidList, clusterUuid.RawBytes())
					continue
				}

				effectiveRateLimit, err := getEffectiveRateLimit(ctx, cpdbIfc, categoryIfc, zeusConfig, clusterUuid,
					"get_ecap_for_rate_limit")
				if err != nil {
					return err
				}

				if effectiveRateLimit == nil {
					errMsg := fmt.Sprintf("Could not resolve effective rate limit for cluster %s", clusterUuid.String())
					return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
				}

				clusterStatus := &marinaIfc.CatalogRateLimitInfo_CatalogRateLimitStatusInfo_ClusterRateLimitStatusInfo{
					ClusterUuid:            clusterUuid.RawBytes(),
					EffectiveRateLimitKbps: effectiveRateLimit.RateLimitKbps,
					EffectiveRateLimitUuid: effectiveRateLimit.Uuid,
				}

				if isRateLimitAuthorized(effectiveRateLimit, authorized) {
					clusterStatus.EffectiveRateLimitName = effectiveRateLimit.Name
				} else {
					clusterStatus.EffectiveRateLimitName = proto.String(redacted)
				}

				status.ClusterStatusList = append(status.ClusterStatusList, clusterStatus)
			}

			rateLimit.Status = status
			return nil
		}
	}
	return nil
}

func getEffectiveRateLimit(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, categoryIfc categoryUtil.CategoryUtilInterface,
	zeusConfig zeus.ConfigCache, peUuid *uuid4.Uuid, queryName string) (*marinaIfc.CatalogRateLimitInfo, error) {

	isSupported := zeusConfig.IsRateLimitSupported(peUuid)
	if isSupported == nil {
		errMsg := "invalid zeus config. IsRateLimitSupported property not found"
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	if !*isSupported {
		log.Infof("Rate limit not supported on %s", peUuid.String())
		return nil, nil
	}

	query, err := QUERY(queryName).
		SELECT(ecapAttributes...).
		FROM(db.EntityCapability.ToString()).
		WHERE(AND(
			EQ(COL(kindCol), STR(clusterKind)),
			EQ(COL(kindIdCol), STR(peUuid.String())))).
		Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create the IDF query %s: %v", queryName, err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	entities, err := cpdbIfc.Query(arg)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch the ecap entry for cluster %s: %v", peUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	if len(entities) != 1 {
		errMsg := fmt.Sprintf("Found more than 1 ecap entry for cluster %s", peUuid.String())
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	categoryIds, err := entities[0].GetStringList(categoryIdListCol)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get category ID List for cluster %s", peUuid.String())
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	var categoryUuids []*uuid4.Uuid
	for _, categoryUuidStr := range categoryIds {
		categoryUuid, err := uuid4.StringToUuid4(categoryUuidStr)
		if err != nil {
			errMsg := fmt.Sprintf("Error encountered while creating category UUID: %v", err)
			return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
		}
		categoryUuids = append(categoryUuids, categoryUuid)
	}

	rateLimitUuids, err := categoryIfc.GetPolicyUuidsByKindForCategoryUuids(ctx, []string{rateLimitFilterKind},
		categoryUuids, true)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get rate limit UUIDs for category IDs %v", categoryIds)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	if len(rateLimitUuids) == 0 {
		log.Infof("No matching rate limits for %s", peUuid.String())
		return nil, nil
	}

	var guids []*insights_interface.EntityGuid
	for _, uuid := range rateLimitUuids {
		guids = append(guids, &insights_interface.EntityGuid{
			EntityTypeName: proto.String(db.RateLimit.ToString()),
			EntityId:       proto.String(uuid.String()),
		})
	}

	rateLimitEntities, err := cpdbIfc.GetEntities(guids, false)
	if err != nil {
		errMsg := "error encountered while getting Rate Limits"
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	var rateLimits []*marinaIfc.CatalogRateLimitInfo
	for _, entity := range rateLimitEntities {
		rateLimit := &marinaIfc.CatalogRateLimitInfo{}
		err = entity.DeserializeEntity(rateLimit)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to deserialize rate limit IDF entry: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		rateLimits = append(rateLimits, rateLimit)
	}

	sort.SliceStable(rateLimits, func(i, j int) bool {
		if rateLimits[i].GetRateLimitKbps() != rateLimits[j].GetRateLimitKbps() {
			return rateLimits[i].GetRateLimitKbps() < rateLimits[j].GetRateLimitKbps()
		}
		return uuid4.ToUuid4(rateLimits[i].GetUuid()).String() < uuid4.ToUuid4(rateLimits[j].GetUuid()).String()
	})

	return rateLimits[0], nil
}

func isRateLimitAuthorized(rateLimit *marinaIfc.CatalogRateLimitInfo, authorized *authz.Authorized) bool {
	if authorized.All {
		return true
	}

	for _, uuid := range authorized.Uuids.ToSlice() {
		if uuid.Equals(uuid4.ToUuid4(rateLimit.GetUuid())) {
			return true
		}
	}

	return false
}
