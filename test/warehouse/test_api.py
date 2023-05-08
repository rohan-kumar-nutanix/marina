import ntnx_marina_py_client as client
from ntnx_marina_py_client.Ntnx.cms.v4.config import (
    SecurityPolicy, ScannerConfig, TrivyServiceConfig)
from ntnx_marina_py_client.Ntnx.cms.v4.content import (
    Warehouse, WarehouseItem, DockerImageReference)

from ntnx_marina_py_client.Ntnx.cms.v4.config import *
from  ntnx_marina_py_client.Ntnx.cms.v4.content import ItemType
c = client.configuration.Configuration()
c.username = "admin"
c.password = "Nutanix.123"
my_pc = "10.37.160.77"
c.debug = True
c.verify_ssl = False
c.host = "10.37.160.77"

api_client = client.ApiClient(c)
api_client.add_default_header(
    header_name="Authorization",
    header_value=c.get_basic_auth_token()
)

# api_instance = client.WarehouseApi(api_client)
warehouse_api =  client.WarehouseApi(api_client)
warehouse_item_api = client.WarehouseItemsApi(api_client)
sp_api = client.SecurityPolicyApi(api_client)
scanner_api = client.ScannerToolsApi(api_client)


def CreateWarehouse():
    w = Warehouse.Warehouse()
    w.name = "Harbour Repo"
    w.description = "Container Images Test desc"

    res = warehouse_api.create_warehouse(w)
    print("response %s", res)


def AddItemToWarehouse():
    item = WarehouseItem.WarehouseItem()
    item.name = "Nginx Image"
    item.description = "V1 Nginx small size"
    item.type = ItemType.ItemType.CONTAINER_IMAGE
    docker_img = DockerImageReference.DockerImageReference(sha256_hex_digest="491b9802c21c")
    item.source_reference = docker_img
    print("WarehouseItem Create",
          warehouse_item_api.add_item_to_warehouse(extId="f10a242a-02f2-46ff-7477-bb6e2b903580", body=item))


def ListWarehouses():
    print("List of Warehouses :\n", warehouse_api.list_warehouses())

def GetWarehouseByExtID():
    print("List of Warehouses :\n", warehouse_api.get_warehouse(extId="f10a242a-02f2-46ff-7477-bb6e2b903580"))

def UpdateWarehouse():
    w = Warehouse.Warehouse()
    w.name = "Update Harbour Repo"
    w.description = "Update Container Images Test desc"
    res = warehouse_api.update_warehouse_metadata(extId="f10a242a-02f2-46ff-7477-bb6e2b903580", body=w)
    print("response %s", res)


def CreateSecurityPolicy():
    sp = SecurityPolicy.SecurityPolicy()
    sp.name = "Daily Scanner"
    sp.description = "Scan the Image Daily"
    sp.image_types = [ItemType.ItemType.DISK_IMAGE]
    print("Create Security Policy Response :\n",sp_api.create_security_policy(sp))


def CreateScannerToolConfiguration():
    config = ScannerConfig.ScannerConfig()
    config.name = "Docker Swarm"
    config.description = "Image Scanner to scan Docker Images"
    config.image_types = [ItemType.ItemType.CONTAINER_IMAGE, ItemType.ItemType.DISK_IMAGE]
    config.server_type = ScannerConfig.ScanServerType.TRIVY

    trivyConfig = TrivyServiceConfig.TrivyServiceConfig()
    trivyConfig.service_url = "https://trivy_server:9090/service"
    trivyConfig.access_token = "asdf_token"
    trivyConfig.is_secure = True

    config.server_config = trivyConfig

    print("Create ScannerTool Configuration :\n", scanner_api.create_scanner_tool_config(config))

#### Warehouse and WarehouseItem WorkFlows
# CreateWarehouse()
# AddItemToWarehouse()
# UpdateWarehouse()
# UpdateWarehouseItem()

# ListWarehouses()
# GetWarehouseByExtID()

#### SecurityPolicy Workflows
# CreateSecurityPolicy()


#### ScannerTools Configuration Workflows
CreateScannerToolConfiguration()
