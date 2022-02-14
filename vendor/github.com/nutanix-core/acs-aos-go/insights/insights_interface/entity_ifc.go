/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Defines Entity interface.
 */

package insights_interface

type EntityIfc interface {
  GetString(attrName string) (string, error)
  GetStringList(attrName string) ([]string, error)
}
