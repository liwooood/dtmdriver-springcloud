package driver

import (
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"testing"
)

func TestNacosDriver_RegisterHttpService(t *testing.T) {
	target := "127.0.0.1:8848"
	endpoint := "127.0.0.1:36789"

	options := make(map[string]string)
	options["username"] = "nacos"
	options["password"] = "nacos"

	driver := &nacosDriver{}
	err := driver.RegisterHttpService(target, endpoint, options)
	if err != nil {
		logger.Errorf("exceprion: %v", err)
		t.FailNow()
	}

	//time.Sleep(3600 * time.Second)
}

func TestNacosDriver_SelectOneHealthyInstancer(t *testing.T) {
	target := "127.0.0.1:8848"
	endpoint := "127.0.0.1:36789"

	options := make(map[string]string)
	options["username"] = "nacos"
	options["password"] = "nacos"

	driver := &nacosDriver{}
	err := driver.RegisterHttpService(target, endpoint, options)
	if err != nil {
		logger.Errorf("exception: %v", err)
		t.Failed()
	}

	instance, err := driver.SelectOneHealthyInstance("dtmService", "", nil)
	if err != nil {
		logger.Errorf("exception: %v", err)
		t.Failed()
	}
	logger.Infof("%v", instance)
}
