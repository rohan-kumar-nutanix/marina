/*
 * Copyrigright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Adds helper methods to 'Entity' protobuf generated object.
 */

package insights_interface

type Vm struct {
  entity EntityIfc
}

func NewVm(entity EntityIfc) *Vm {
  return &Vm{
    entity: entity,
  }
}

//-----------------------------------------------------------------------------

func (vm *Vm) GetIPs() ([]string, error) {
  ips, err := vm.entity.GetStringList("ip_addresses")
  if err != nil {
    return nil, err
  }
  return ips, nil
}
