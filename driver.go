package driver

import (
	"errors"
	"github.com/dtm-labs/dtmdriver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

const (
	DriverName  = "dtm-driver-nacos"
	ServiceName = "serviceName"
)

type nacosDriver struct {
	nacosClient naming_client.INamingClient
}

func (n *nacosDriver) GetName() string {
	return DriverName
}

func (n *nacosDriver) RegisterGrpcResolver() {
}

func (n *nacosDriver) RegisterGrpcService(target string, endpoint string) error {
	//TODO implement me
	panic("implement me")
}

func (n *nacosDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	//TODO implement me
	panic("implement me")
}

func (n *nacosDriver) RegisterHttpResolver() {

}

func (n *nacosDriver) RegisterHttpService(target string, endpoint string, options map[string]string, paths []string) error {
	if n.nacosClient == nil {
		err := n.buildNacosClient(target, options)
		if err != nil {
			return err
		}
	}
	ip, port, err := splitIpAndPort(endpoint)
	if err != nil {
		return err
	}

	registerParam := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: ServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	if v, ok := options["clusterName"]; ok {
		registerParam.ClusterName = v
	}
	succ, err := n.nacosClient.RegisterInstance(registerParam)
	if err != nil {
		return err
	}
	if !succ {
		logger.Infof("register service %s to nacos fail.", ServiceName)
	}
	for _, path := range paths {
		logger.Infof("current service name is : %s", path)
		//succ, err := n.nacosClient.RegisterInstance(vo.RegisterInstanceParam{
		//	Ip:          ip,
		//	Port:        port,
		//	ServiceName: ServiceName,
		//	Weight:      10,
		//	ClusterName: DtmService,
		//	Enable:      true,
		//	Healthy:     true,
		//	Ephemeral:   true,
		//})
		//if err != nil {
		//	return err
		//}
		//if !succ {
		//	logger.Infof("register service %s to nacos fail.", path)
		//}
	}
	return nil
}

func (n *nacosDriver) buildNacosClient(target string, options map[string]string) error {
	ip, port, err := splitIpAndPort(target)
	if err != nil {
		return err
	}

	if _, ok := options["username"]; !ok {
		return errors.New("must configure username for nacos")
	}
	if _, ok := options["password"]; !ok {
		return errors.New("must configure password for nacos")
	}

	namespaceId := ""
	if v, ok := options["namespaceId"]; ok {
		namespaceId = v
	}
	var timeout uint64 = 5000
	if v, ok := options["timeout"]; ok {
		num, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		timeout = num
	}
	var listenInterval uint64 = 10000
	if v, ok := options["listenInterval"]; ok {
		num, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		listenInterval = num
	}
	notLoadCacheAtStart := true
	if v, ok := options["notLoadCacheAtStart"]; ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		notLoadCacheAtStart = b
	}
	logDir := "data\\nacos\\log"
	if v, ok := options["logDir"]; ok {
		logDir = v
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{
			{
				IpAddr: ip,
				Port:   port,
			},
		},
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           timeout,
			ListenInterval:      listenInterval,
			NotLoadCacheAtStart: notLoadCacheAtStart,
			LogDir:              logDir,
			Username:            options["username"],
			Password:            options["password"],
			NamespaceId:         namespaceId,
		},
	})

	n.nacosClient = client
	return nil
}

func splitIpAndPort(str string) (string, uint64, error) {
	ipPort := strings.Split(str, ":")
	ip := ipPort[0]
	port, err := strconv.ParseUint(ipPort[1], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}

func init() {
	dtmdriver.Register(&nacosDriver{})
}
