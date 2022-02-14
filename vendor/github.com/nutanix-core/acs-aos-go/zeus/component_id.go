package zeus

import (
	"errors"
	"fmt"
	"github.com/nutanix-core/acs-aos-go/go-zookeeper"
	"regexp"
	"strconv"
	"time"
)

const zookeeper_component_id_path = "/appliance/logical/variables/component-id"

func Allocate_component_ids(conn *zk.Conn, num_ids int) (int64, error) {
	if num_ids <= 0 {
		fmt.Sprintf("Invalid number of ids requested: %d", num_ids)
		return -1, errors.New("Invalid number of ids requested")
	}
	for {
		data, stat, err := conn.GetWithBackoff(zookeeper_component_id_path)
		// TODO(Hitesh): In case of zk.ErrNoNode, create the node.
		if err != nil {
			fmt.Printf("Failed to retrieve component-id zknode with error: %s", err)
			return -1, errors.New("Failed to retrieve component-id zknode")
		}
		r, _ := regexp.Compile("value: (\\d+)")
		match := r.FindStringSubmatch(string(data))
		if len(match) != 2 {
			fmt.Printf("Failed to parse component-id zknode: %s", string(data))
			return -1, errors.New("Failed to parse component-id zknode")
		}

		allocated_id, err := strconv.ParseInt(match[1], 10 /* base */, 64 /* bitSize */)
		val := "value: " + strconv.FormatInt(allocated_id + int64(num_ids), 10 /* base */)
		stat, err = conn.SetWithBackoff(zookeeper_component_id_path, []byte(val), stat.Version)
		if err == nil {
			return allocated_id, nil
		}
		time.Sleep(time.Second)
	}
}
