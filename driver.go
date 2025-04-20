package driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dtm-labs/dtmdriver"
	"github.com/go-resty/resty/v2"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type config struct {
	Addr           string
	Type           string
	ClientConfig   constant.ClientConfig
	InstanceConfig vo.RegisterInstanceParam
}

type springcloudDriver struct {
	client naming_client.INamingClient
}

func (d *springcloudDriver) GetName() string {
	return "dtm-driver-springcloud"
}

func (d *springcloudDriver) RegisterService(target string, endpoint string) error {
	conf := config{}
	err := json.Unmarshal([]byte(target), &conf)
	if err != nil {
		return fmt.Errorf("invalid options: %s error: %w", target, err)
	}
	if conf.Type == "nacos" {
		d.client, err = newNacosClient(strings.Split(conf.Addr, ","), conf.ClientConfig)
		if err != nil {
			return fmt.Errorf("new nacos client error: %w", err)
		}
	} else {
		return fmt.Errorf("unknown type: %s", conf.Type)
	}
	return d.registerService(conf.InstanceConfig, endpoint)
}

func (d *springcloudDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	return "", "", nil
}

func (d *springcloudDriver) RegisterAddrResolver() {
	dtmdriver.Middlewares.HTTP = append(dtmdriver.Middlewares.HTTP, func(c *resty.Client, r *resty.Request) (err error) {
		r.URL, err = d.resolveURL(r.URL)
		return
	})
}

func init() {
	dtmdriver.Register(&springcloudDriver{})
}
