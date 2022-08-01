package driver

import (
	nurl "net/url"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func newNacosClient(addrs []string, clientConfig constant.ClientConfig) (naming_client.INamingClient, error) {
	serverConfigs := []constant.ServerConfig{}
	for _, addr := range addrs {
		ip, port, err := parseIpPort(addr)
		if err != nil {
			return nil, err
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip,
			Port:   port,
		})
	}

	nclient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return nclient, nil
}

func (d *springcloudDriver) resolveURL(url string) (string, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url, nil
	}
	u, err := nurl.Parse(url)
	if err != nil {
		return "", err
	}
	groupName := u.Query().Get("group")
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	clusters := u.Query()["clusters"]
	if len(clusters) == 0 {
		clusters = append(clusters, "DEFAULT")
	}
	instance, err := d.client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: u.Host,
		GroupName:   groupName,
		Clusters:    clusters,
	})
	if err != nil || instance == nil {
		return url, err
	}
	u.Host = instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
	u.Scheme = "http"
	return u.String(), nil
}

func (d *springcloudDriver) registerService(config vo.RegisterInstanceParam, endpoint string) error {
	ip, port, err := parseIpPort(endpoint)
	if err != nil {
		return err
	}
	config.Ip = ip
	config.Port = port
	if config.Weight == 0 {
		config.Weight = 10
	}
	logger.Infof("registering service: %v", config)
	_, err = d.client.RegisterInstance(config)
	return err
}
