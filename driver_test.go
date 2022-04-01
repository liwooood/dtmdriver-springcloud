package driver

import (
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"testing"
)

func TestNacosDriver_RegisterHttpService(t *testing.T) {
	target := "121.4.131.37:8848"
	endpoint := "127.0.0.1:26789"

	options := make(map[string]string)
	options["username"] = "nacos"
	options["password"] = "nacos"

	paths := make([]string, 0)
	paths = append(paths, "/api/dtmsvr/newGid")
	paths = append(paths, "/api/dtmsvr/prepare")
	paths = append(paths, "/api/dtmsvr/submit")

	driver := &nacosDriver{}
	err := driver.RegisterHttpService(target, endpoint, options, paths)
	if err != nil {
		logger.Errorf("exceprion: %v", err)
		t.FailNow()
	}

	//time.Sleep(3600 * time.Second)

}
