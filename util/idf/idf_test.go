/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * IDF utilities for Marina Service.
 */

package idf

import (
	"context"
	"reflect"
	"testing"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
)

func TestIdfClient_Query(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		IdfSvc insights_interface.InsightsServiceInterface
		Retry  *misc.ExponentialBackoff
	}
	type args struct {
		arg *insights_interface.GetEntitiesWithMetricsArg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*insights_interface.EntityWithMetric
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idf := &IdfClient{
				IdfSvc: tt.fields.IdfSvc,
				Retry:  tt.fields.Retry,
			}
			got, err := idf.Query(ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("IdfClient.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IdfClient.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
