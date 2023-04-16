API_PROTO_PATH=./protos/apis
CMS_CONFIG_PROTO_PATH=$API_PROTO_PATH/cms/v4/config
CMS_CONTENT_PROTO_PATH=$API_PROTO_PATH/cms/v4/content
CMS_ERROR_PROTO_PATH=$API_PROTO_PATH/cms/v4/error
COMMON_RESPONSE_PROTO=$API_PROTO_PATH/common/v1/response
SED_OPTS=" -i -e "

#if [[ "$OSTYPE" == "darwin"* ]]; then
#SED_OPTS="-i \'\' -e"
#fi

if [[ "$OSTYPE" == "darwin"* ]]; then
#content.pb.go
sed -i '' -e 's/cms\/v4\/error"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/error"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed -i '' -e 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed -i '' -e 's/common\/v1\/response"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/response"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed -i '' -e 's/prism\/v4\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/prism\/v4\/config"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go

#config.pb.go
sed -i '' -e 's/cms\/v4\/error"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/error"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed -i '' -e 's/cms\/v4\/content"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/content"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed -i '' -e 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed -i '' -e 's/common\/v1\/response"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/response"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed -i '' -e 's/prism\/v4\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/prism\/v4\/config"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go

#response.pb.go
sed -i '' -e 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $COMMON_RESPONSE_PROTO/response.pb.go

# ScannerTools_service.pb.go
sed -i '' -e 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONFIG_PROTO_PATH/ScannerTools_service.pb.go

# protos/apis/cms/v4/config/SecurityPolicy_service.proto
sed -i '' -e 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONFIG_PROTO_PATH/SecurityPolicy_service.pb.go

# protos/apis/cms/v4/content/Warehouse_service.proto
sed -i '' -e 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONTENT_PROTO_PATH/Warehouse_service.pb.go

#protos/apis/cms/v4/content/WarehouseItems_service.proto
sed -i '' -e 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONTENT_PROTO_PATH/WarehouseItems_service.pb.go

#protos/apis/cms/v4/error/error.proto
sed -i '' -e 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_ERROR_PROTO_PATH/error.pb.go

else
#content.pb.go
sed $SED_OPTS 's/cms\/v4\/error"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/error"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed $SED_OPTS 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed $SED_OPTS 's/common\/v1\/response"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/response"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go
sed $SED_OPTS 's/prism\/v4\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/prism\/v4\/config"/g' $CMS_CONTENT_PROTO_PATH/content.pb.go

#config.pb.go
sed $SED_OPTS 's/cms\/v4\/error"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/error"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed $SED_OPTS 's/cms\/v4\/content"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4\/content"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed $SED_OPTS 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed $SED_OPTS 's/common\/v1\/response"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/response"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go
sed $SED_OPTS 's/prism\/v4\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/prism\/v4\/config"/g' $CMS_CONFIG_PROTO_PATH/config.pb.go

#response.pb.go
sed $SED_OPTS 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $COMMON_RESPONSE_PROTO/response.pb.go

# ScannerTools_service.pb.go
sed $SED_OPTS 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONFIG_PROTO_PATH/ScannerTools_service.pb.go

# protos/apis/cms/v4/config/SecurityPolicy_service.proto
sed $SED_OPTS 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONFIG_PROTO_PATH/SecurityPolicy_service.pb.go

# protos/apis/cms/v4/content/Warehouse_service.proto
sed $SED_OPTS 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONTENT_PROTO_PATH/Warehouse_service.pb.go

#protos/apis/cms/v4/content/WarehouseItems_service.proto
sed $SED_OPTS 's/cms\/v4"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/cms\/v4"/g' $CMS_CONTENT_PROTO_PATH/WarehouseItems_service.pb.go

#protos/apis/cms/v4/error/error.proto
sed $SED_OPTS 's/common\/v1\/config"/github.com\/nutanix-core\/content-management-marina\/protos\/apis\/common\/v1\/config"/g' $CMS_ERROR_PROTO_PATH/error.pb.go

fi
