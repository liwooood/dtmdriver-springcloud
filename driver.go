package driver

import (
	"fmt"
	"strings"

	"github.com/dtm-labs/dtmdriver"
)

type httpDriver struct {
	client dtmdriver.HTTPClient
}

func (d *httpDriver) GetName() string {
	return "dtm-driver-http"
}

func (d *httpDriver) ResolveURL(url string) (string, error) {
	return d.client.ResolveURL(url)
}

func (d *httpDriver) RegisterService(target string, endpoint string) error {
	return d.client.RegisterService(target, endpoint)
}

func (d *httpDriver) Init(registryType string, address string, options string) error {
	addrs := strings.Split(address, ",")
	var err error
	if registryType == "nacos" {
		d.client, err = newNacosClient(addrs, options)
	} else {
		return fmt.Errorf("unsupported registry type: %s", registryType)
	}
	return err
}
