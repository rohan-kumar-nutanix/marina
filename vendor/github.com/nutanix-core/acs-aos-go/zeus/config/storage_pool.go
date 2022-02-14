/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Add helper routines to proto generated 'ConfigurationProto_StorageTier'
 *
 */

package zeus_config

type ByRandomPriority []*ConfigurationProto_StorageTier

func (a ByRandomPriority) Len() int      { return len(a) }
func (a ByRandomPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRandomPriority) Less(i, j int) bool {
        return *a[i].RandomIoPriority < *a[j].RandomIoPriority
}

type BySequentialPriority []*ConfigurationProto_StorageTier

func (a BySequentialPriority) Len() int      { return len(a) }
func (a BySequentialPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySequentialPriority) Less(i, j int) bool {
        return *a[i].SequentialIoPriority < *a[j].SequentialIoPriority
}
