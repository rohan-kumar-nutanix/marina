package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/glog"

	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
)

// Flags.
var (
	factoryConfigPath = flag.String(
		"test_only_factory_config_path",
		"/etc/nutanix/factory_config.json",
		"This value is overwritten only for integration tests.")
)

const (
	FACTORY_CONFIG_NODE_UUID_IDENT = "node_uuid"
)

var (
	testFactoryConfigJson = []byte(`
	{
		"rackable_unit_serial": "10-5-22-148",
		"node_uuid": "f81bfc7b-377c-440c-927d-d20cdf31e50a",
		"node_serial": "10-5-22-148",
		"node_position": "A",
		"cluster_id": 999,
		"rackable_unit_model": "null"
	}
	`)
)

// Config Error definition.
type configError struct {
	*ntnx_errors.NtnxError
}

func (e *configError) TypeOfError() int {
	return ntnx_errors.ConfigErrorType
}

func ConfigError(errStr string, errCode int) *configError {
	return &configError{ntnx_errors.NewNtnxError(errStr, errCode)}
}
func (e *configError) SetCause(err error) *configError {
	return &configError{e.NtnxError.SetCause(err).(*ntnx_errors.NtnxError)}
}

var (
	ErrFileRead   = ConfigError("File read error", 1)
	ErrUnmarshall = ConfigError("Unmarshalling error", 2)
)

func GetUuidFromFactoryConfig() (*uuid4.Uuid, error) {
	factory_contents, err := ioutil.ReadFile(*factoryConfigPath)
	if err != nil {
		glog.Error("Failed reading file at", *factoryConfigPath)
		return nil, ErrFileRead
	}
	var json_data interface{}
	err = json.Unmarshal(factory_contents, &json_data)
	if err != nil {
		glog.Error("Failed unmarshal file content as json at %s",
			*factoryConfigPath)
		return nil, ErrUnmarshall
	}
	json := json_data.(map[string]interface{})
	uuid := json[FACTORY_CONFIG_NODE_UUID_IDENT].(string)
	return uuid4.StringToUuid4(uuid)
}

// Setup a mock factory config JSON. Expected to be run within a unit test,
// where the $PWD would be something like "/tmp/go..."
func TestSetupMockFactoryConfig() {
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	mockFactoryConfig := filepath.Join(pwd, "mock_factory_config.json")
	flag.Set("test_only_factory_config_path", mockFactoryConfig)

	if err := ioutil.WriteFile(mockFactoryConfig, testFactoryConfigJson,
		os.ModePerm); err != nil {
		panic(err)
	}
}
