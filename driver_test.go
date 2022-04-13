package driver

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNacosDriver_RegisterHttpService(t *testing.T) {
	target := "1.14.127.29:8848"
	endpoint := "127.0.0.1:36789"

	options := make(map[string]string)
	options["username"] = "nacos"
	options["password"] = "nacos"

	driver := &nacosDriver{}
	err := driver.RegisterHttpService(target, endpoint, options)
	assert.Nil(t, err)

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
	assert.Nil(t, err)

	instance, err := driver.SelectOneHealthyInstance("dtmService", "", nil)
	assert.Nil(t, err)
	logger.Infof("%v", instance)
}

func TestNacosDriver_ResolveHttpService(t *testing.T) {
	target := "127.0.0.1:8848"
	endpoint := "127.0.0.1:36789"

	options := make(map[string]string)
	options["username"] = "nacos"
	options["password"] = "nacos"
	options["namespaceId"] = "c3dc917d-906a-429d-90a9-85012b41014e"

	driver := &nacosDriver{}
	err := driver.RegisterHttpService(target, endpoint, options)
	if err != nil {
		logger.Errorf("exception: %v", err)
		t.Failed()
	}

	originalUrl := "http://dtmcli-spring-boot-starter-sample/feignTest?clusters=[\"DEFAULT\"]&groupName=DEFAULT_GROUP&a=1"
	realUrl := driver.ResolveHttpService(originalUrl)
	assert.Equal(t, realUrl, "http://192.168.101.9:8888/feignTest?a=1")
	fmt.Println(realUrl)

	originalUrl = "http://192.168.101.9:8888/feignTest?a=1"
	realUrl = driver.ResolveHttpService(originalUrl)
	assert.Equal(t, realUrl, originalUrl)

	originalUrl = "http://dtmcli-spring-boot-starter-sample/feignTest?groupName=abcd"
	realUrl = driver.ResolveHttpService(originalUrl)
	assert.Equal(t, realUrl, originalUrl)
}
