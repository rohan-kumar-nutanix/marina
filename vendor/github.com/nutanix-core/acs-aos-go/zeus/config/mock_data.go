/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Stateless functions around 'zeus/config'. Very convenient, particularly if
 * you don't intend to save the configuration proto, zksession, etc. But
 * expensive b/c getting a zksession each time is expensive.
 */

package zeus_config

import (
        "github.com/golang/protobuf/proto"
)

func MockConfig1() *ConfigurationProto {
        config := &ConfigurationProto{}
        proto.UnmarshalText(`
logical_timestamp: 34
internal_subnet: "10.5.20.0/255.255.252.0"
external_subnet: "10.5.20.0/255.255.252.0"
storage_tier_list: <
  storage_tier_name: "SSD-PCIe"
  random_io_priority: 9
  sequential_io_priority: 9
  storage_type: kPcieSSD
>
storage_tier_list: <
  storage_tier_name: "SSD-SATA"
  random_io_priority: 8
  sequential_io_priority: 8
  storage_type: kSataSSD
>
storage_tier_list: <
  storage_tier_name: "DAS-SATA"
  random_io_priority: 7
  sequential_io_priority: 7
  storage_type: kHDD
>
disk_list: <
  disk_id: 9
  service_vm_id: 2
  mount_path: "/home/nutanix/data/stargate-storage/disks/NFS_4_0_289_7c174117_320b_480f_ac05_b0f622127511"
  disk_size: 464566805299
  statfs_disk_size: 522925326336
  storage_tier: "SSD-SATA"
  data_dir_sublevels: 2
  data_dir_sublevel_dirs: 20
  data_migration_status: 0
  contains_metadata: true
  oplog_disk_size: -1
  disk_serial_id: "NFS_4_0_289_7c174117_320b_480f_ac05_b0f622127511"
  disk_uuid: "c27bf065-2f53-43e0-a480-f2c37d1e32c8"
  node_uuid: "f81bfc7b-377c-440c-927d-d20cdf31e50a"
  chosen_for_metadata: false
  metadata_disk_reservation_bytes: 32212254720
>
node_list: <
  service_vm_id: 2
  service_vm_external_ip: "10.5.22.148"
  node_status: kNormal
  cassandra_token_id: "0000000042c1veYNsUO5xGqZ0GP85xeVovcCN1rwlUCBCGw85zIyBuepC9Hc"
  hypervisor: <
    username: "root"
  >
  zookeeper_myid: 1
  uuid: "f81bfc7b-377c-440c-927d-d20cdf31e50a"
  rackable_unit_id: 8
  node_position: 1
  software_version: "el7.3-opt-master-c4ba1413b4af016671cb762d737f74e49b6920e6"
  node_serial: "10-5-22-148"
  cluster_uuid: "cfa38612-1ace-4235-8815-263462867b2d"
  rackable_unit_uuid: "af630aed-e611-498f-b5b3-1099e01b8344"
  cassandra_schema_timestamp: 2
  controller_vm_backplane_ip: "10.5.22.148"
>
storage_pool_list: <
  storage_pool_name: "default-storage-pool-11234176403077"
  storage_pool_id: 3
  disk_id: 9
  storage_pool_uuid: "012dfa6d-657b-4019-8d7a-6a00add7ed75"
  disk_uuid: "c27bf065-2f53-43e0-a480-f2c37d1e32c8"
>
container_list: <
  container_name: "default-container-11234176403077"
  container_id: 4
  storage_pool_id: 3
  params: <
    replication_factor: 1
    random_io_tier_preference: "SSD-PCIe"
    random_io_tier_preference: "SSD-SATA"
    random_io_tier_preference: "DAS-SATA"
    sequential_io_tier_preference: "SSD-PCIe"
    sequential_io_tier_preference: "SSD-SATA"
    sequential_io_tier_preference: "DAS-SATA"
    max_capacity: 464566805299
    ilm_down_migrate_time_secs: 1800
    ilm_down_migrate_time_secs: 1800
    ilm_down_migrate_time_secs: 1800
    oplog_params: <
      replication_factor: 1
      num_stripes: 1
      need_sync: false
    >
    fingerprint_on_write: false
  >
  container_uuid: "764bc38b-4c44-4877-8920-25dd624867a6"
  storage_pool_uuid: "012dfa6d-657b-4019-8d7a-6a00add7ed75"
>
container_list: <
  container_name: "NutanixManagementShare"
  container_id: 5
  storage_pool_id: 3
  params: <
    replication_factor: 1
    random_io_tier_preference: "SSD-PCIe"
    random_io_tier_preference: "SSD-SATA"
    random_io_tier_preference: "DAS-SATA"
    sequential_io_tier_preference: "SSD-PCIe"
    sequential_io_tier_preference: "SSD-SATA"
    sequential_io_tier_preference: "DAS-SATA"
    max_capacity: 464566805299
    ilm_down_migrate_time_secs: 1800
    ilm_down_migrate_time_secs: 1800
    ilm_down_migrate_time_secs: 1800
    oplog_params: <
      replication_factor: 1
      num_stripes: 1
      need_sync: false
    >
    fingerprint_on_write: false
    nutanix_managed: true
  >
  container_uuid: "f2a1d079-6688-41cd-81b9-72ed9c0f2809"
  storage_pool_uuid: "012dfa6d-657b-4019-8d7a-6a00add7ed75"
>
cluster_incarnation_id: 5738577762517402165
cluster_id: 582413733247482669
aegis: <
  remote_support: <
    value: false
    until_time_usecs: 0
  >
  email_alerts: <
    value: true
    until_time_usecs: 0
  >
  auto_support_config: <
    email_asups: <
      value: false
      until_time_usecs: 0
    >
    send_email_asups_externally: false
    aos_version: "master"
    last_login_workflow_time_msecs: 1500318586157
    remind_later: false
  >
  send_email_alerts_externally: false
  send_alert_email_digest: true
>
default_gateway_ip: "10.5.20.1"
rackable_unit_list: <
  rackable_unit_id: 8
  rackable_unit_serial: "null"
  rackable_unit_model: kNull
  rackable_unit_model_name: "null"
  rackable_unit_uuid: "af630aed-e611-498f-b5b3-1099e01b8344"
>
cassandra_schema_version: "el7.3-opt-master-c4ba1413b4af016671cb762d737f74e49b6920e6"
vstore_list: <
  vstore_id: 4
  vstore_name: "default-container-11234176403077"
  container_id: 4
  vstore_uuid: "764bc38b-4c44-4877-8920-25dd624867a6"
  container_uuid: "764bc38b-4c44-4877-8920-25dd624867a6"
>
vstore_list: <
  vstore_id: 5
  vstore_name: "NutanixManagementShare"
  container_id: 5
  vstore_uuid: "f2a1d079-6688-41cd-81b9-72ed9c0f2809"
  container_uuid: "f2a1d079-6688-41cd-81b9-72ed9c0f2809"
>
release_version: "el7.3-opt-master-c4ba1413b4af016671cb762d737f74e49b6920e6"
ssh_key_list: <
  key_id: "agave"
  pub_key: "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAw8NhpiLG/6c8c0LocaQVM6+x8yBvLS+z13nYlqzSKr64HQSH+IIx14bS+uh0xf8ou5ROnbUa44YM9qG2RBK1hBVrwTIKZBwFwOTxDuXH9YoswqflnOeyTGUFerrMxlg5iYvKE01QIErKWt/BJ+hK/qVQi6SHh/WhyjAGNQPxjPBG3oEUFoLV9uDLX9RT/sLIQpBf91DCbF9+2rSToqjydj36d38Jbmu
GJGqr+DrEd0lKzfiJiQ39t3AXBE1bY2ISX+uwq07XkYZbrHaijrxnOBFZE7ETQmyxdcOvgDxvhnhZpFsUqbPqIqa45CwFWV+btzITvpJkveYjGSGR7ZJNHQ== nutanix@titan"
  key_type: kExternalKey
  ssl_cert: ""
>
ssh_key_list: <
  key_id: "f81bfc7b-377c-440c-927d-d20cdf31e50a"
  pub_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDf5sguzg6w7OvxRF/NRG0jiOtfsMp7lrT2AY4sHT1EgGioQ8DfvZjr79bPNOFwjrqjJ8dqM+J9pxPPFLnEtjStUcCyfK2+vLQ4xrJrFT31c64tI0M0YMoy8lUGI4n9le0nLOsUhv6oTu1E0hxQutwyQShzZugj8/mVJo2rkT7EELDdqa0+r5Y3D4KGPOVKyEMiSDYeqHRP8lvLEuShhLBOrOjmvKb
aiS/w9iy8sk69QzytSuiUQNW9yZCuVs+9W47/1SUgw/K91Geu3G7hmwCp/AKoTNSS2ZCozS8gWwrwIV3zZx+Aco9nfaw4Cy/DjxVs7MqGoL/U6FqFSulkGqD5 nutanix@NTNX-10-5-22-148-A-CVM"
  key_type: kNodeKey
  ssl_cert: "-----BEGIN CERTIFICATE-----\nMIIDrzCCApegAwIBAgIJALWGYDqttipTMA0GCSqGSIb3DQEBCwUAMG4xCzAJBgNV\nBAYTAlVTMQswCQYDVQQIDAJDQTERMA8GA1UEBwwIU2FuIEpvc2UxFTATBgNVBAoM\nDE51dGFuaXggSW5jLjEWMBQGA1UECwwNTWFuYWdlYWJpbGl0eTEQMA4GA1UEAwwH\nbnV0YW5peDAeFw0xNzA3MTcxODQ2MjF
aFw0yNzA3MTUxODQ2MjFaMG4xCzAJBgNV\nBAYTAlVTMQswCQYDVQQIDAJDQTERMA8GA1UEBwwIU2FuIEpvc2UxFTATBgNVBAoM\nDE51dGFuaXggSW5jLjEWMBQGA1UECwwNTWFuYWdlYWJpbGl0eTEQMA4GA1UEAwwH\nbnV0YW5peDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN/myC7ODrDs\n6/FEX81EbSOI61+wynuWtPYBjiwdPUSAaKhDwN
+9mOvv1s804XCOuqMnx2oz4n2n\nE88UucS2NK1RwLJ8rb68tDjGsmsVPfVzri0jQzRgyjLyVQYjif2V7Scs6xSG/qhO\n7UTSHFC63DJBKHNm6CPz+ZUmjauRPsQQsN2prT6vljcPgoY85UrIQyJINh6odE/y\nW8sS5KGEsE6s6Oa8ptqJL/D2LLyyTr1DPK1K6JRA1b3JkK5Wz71bjv/VJSDD8r3U\nZ67cbuGbAKn8AqhM1JLZkKjNLyBbCvAhXfNnH4Byj2d9r
DgLL8OPFWzsyoagv9To\nWoVK6WQaoPkCAwEAAaNQME4wHQYDVR0OBBYEFN3XXVp/RJSQhPBytA+RNiCMj8mE\nMB8GA1UdIwQYMBaAFN3XXVp/RJSQhPBytA+RNiCMj8mEMAwGA1UdEwQFMAMBAf8w\nDQYJKoZIhvcNAQELBQADggEBAFo9NwP8c56hcKxJd192FFCb5EYBOD8Gl8bNRK42\nHiQmwzIt83QXWbxJgEWQwncthSOUDv49QbGSyaWlXpV9ikRZ/20Z
LjVIPrmU8h1s\nZvLPOfFooKhqq0rbjwH2lFFBFKo1nC8yy819OJBplzuzgZjywZALPAHzG0DuCsG0\nZ6TtTbTUmhhRJ9q0h+AUPZV37yX+KHfviJpTsLC31jn1dcK+xEKI5mRSV3jYw9Jz\nZP+mcx37jkqvHc/6WE/a5JWhztG18z7W90wVops9sWO3vww7UPnza/98KpIKG0+x\n1hZ5mVqnb+jWmb7TVdyKhY31QOSbIqpcn8LcGpz6QNhhg6o=\n-----END
CERTIFICATE-----"
>
cluster_fault_tolerance_state: <
  current_max_fault_tolerance: 0
  desired_max_fault_tolerance: 0
>
domain_fault_tolerance_state: <
  domains: <
    domain_type: kDEPRECATEDNode
    components: <
      component_type: kZookeeperInstances
      max_faults_tolerated: 0
      last_update_secs: 1500318363
      tolerance_details_message: <
        message_id: "Zookeeper can tolerate 0 node failure(s)"
      >
    >
    components: <
      component_type: kCassandraRing
      max_faults_tolerated: 0
      last_update_secs: 1500318434
      tolerance_details_message: <
        message_id: "Metadata ring partitions with nodes: {non_fault_tolerant_nodes} are not fault tolerant."
        attribute_list: <
          attribute: "non_fault_tolerant_nodes"
          value: "10.5.22.148"
        >
      >
    >
  >
  domains: <
    domain_type: kDEPRECATEDRackableUnit
    components: <
      component_type: kZookeeperInstances
      max_faults_tolerated: 0
      last_update_secs: 1500318363
      tolerance_details_message: <
        message_id: "Zookeeper can tolerate 0 rackable unit (block) failure(s)"
      >
    >
    components: <
      component_type: kCassandraRing
      max_faults_tolerated: 0
      last_update_secs: 1500318434
      tolerance_details_message: <
        message_id: "Metadata ring partitions with nodes: {non_fault_tolerant_nodes} are not fault tolerant."
        attribute_list: <
          attribute: "non_fault_tolerant_nodes"
          value: "10.5.22.148"
        >
      >
    >
  >
  domains: <
    domain_type: kDEPRECATEDDisk
    components: <
      component_type: kCassandraRing
      max_faults_tolerated: 0
      last_update_secs: 1500318434
      tolerance_details_message: <
        message_id: "Metadata ring partitions with nodes: {non_fault_tolerant_nodes} are not fault tolerant."
        attribute_list: <
          attribute: "non_fault_tolerant_nodes"
          value: "10.5.22.148"
        >
      >
    >
  >
>
cluster_functions: 2
cluster_uuid: "cfa38612-1ace-4235-8815-263462867b2d"
extended_nfs_fhandle_enabled: true
management_share_container_id: 5
        `, config)
        return config
}

func MockNode() *ConfigurationProto_Node {
        node := &ConfigurationProto_Node{}
        proto.UnmarshalText(`
node_list: <
  service_vm_id: 1
  service_vm_external_ip: "10.5.22.147"
  node_status: kNormal
  cassandra_token_id: "..."
  hypervisor: <
    username: "root"
  >
  zookeeper_myid: 1
  uuid: "f81bfc7b-377c-440c-927d-d21cdf31e50a"
  rackable_unit_id: 7
  node_position: 2
  software_version: "el7.3-opt-master-c4ba1413b4af016671cb762d737f74e49b6920e6"
  node_serial: "10-5-22-147"
  cluster_uuid: "cfa38612-1ace-4235-8815-263462867b2d"
  rackable_unit_uuid: "af630aed-e611-498f-b5b3-1099e01b8344"
  cassandra_schema_timestamp: 2
  controller_vm_backplane_ip: "10.5.22.147"
>
`, node)
        return node
}
