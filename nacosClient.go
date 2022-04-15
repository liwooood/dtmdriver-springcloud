package driver

import (
	"encoding/json"
	"fmt"
	nurl "net/url"
	"strconv"

	"github.com/dtm-labs/dtmdriver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosClient struct {
	nacosClient naming_client.INamingClient
}

func (d *nacosClient) ResolveURL(url string) (string, error) {
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
	instance, err := d.nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: u.Host,
		GroupName:   groupName,
		Clusters:    []string{},
	})
	if err != nil || instance == nil {
		return url, err
	}
	u.Host = instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
	return u.String(), nil
}

func (d *nacosClient) RegisterService(target string, endpoint string) error {
	var regParams vo.RegisterInstanceParam
	err := json.Unmarshal([]byte(target), &regParams)
	if err != nil {
		return fmt.Errorf("invalid target: %s error: %w", target, err)
	}

	ip, port, err := parseIpPort(endpoint)
	if err != nil {
		return err
	}
	regParams.Ip = ip
	regParams.Port = port
	if regParams.Weight == 0 {
		regParams.Weight = 10
	}
	_, err = d.nacosClient.RegisterInstance(regParams)
	return err
}

func newNacosClient(addrs []string, options string) (dtmdriver.HTTPClient, error) {
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
	clientConfig := constant.ClientConfig{}
	err := json.Unmarshal([]byte(options), &clientConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid options: %s error: %w", options, err)
	}

	nclient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return &nacosClient{nacosClient: nclient}, nil
}

func init() {
	dtmdriver.RegisterHttp(&httpDriver{})
}
