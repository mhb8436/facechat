package tools

import (
	"context"
	"encoding/json"
	"facechat/config"

	"github.com/sirupsen/logrus"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/metadata"
)

type ServiceInfo struct {
	InstanceId string
	Name       string
	Version    string
	Address    string
	Metadata   metadata.MD
}

func GetAllServiceFromRepository(key string) (ServiceList []ServiceInfo, err error) {
	etcdConfig := etcdv3.Config{
		Endpoints: []string{config.Conf.Common.CommonEtcd.Host},
	}
	etcdCli, err := etcdv3.New(etcdConfig)
	if err != nil {
		logrus.Infof("get etcdcli err")
	}
	resp, err := etcdCli.Get(context.Background(), key, etcdv3.WithPrefix())
	if err == nil {
		ServiceList = extractAddrs(resp)
		logrus.Info("addrs => ", ServiceList, len(ServiceList))
	}
	return
}

func extractAddrs(resp *etcdv3.GetResponse) []ServiceInfo {
	addrs := []ServiceInfo{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			nodeData := ServiceInfo{}
			err := json.Unmarshal(v, &nodeData)
			if err != nil {
				logrus.Info("Parse node data error:", err)
				continue
			}
			addrs = append(addrs, nodeData)
		}
	}
	return addrs
}

func WatchServiceChanged(key string) (ServiceList []ServiceInfo, err error) {
	etcdConfig := etcdv3.Config{
		Endpoints: []string{config.Conf.Common.CommonEtcd.Host},
	}
	etcdCli, err := etcdv3.New(etcdConfig)
	if err != nil {
		logrus.Infof("get etcdcli err")
	}

	respChan := etcdCli.Watch(context.Background(), key, etcdv3.WithPrefix())
	resp := <-respChan
	ServiceList = extractAddrsFromWatchResp(&resp)
	return
}

func extractAddrsFromWatchResp(resp *etcdv3.WatchResponse) []ServiceInfo {
	addrs := []ServiceInfo{}

	if resp == nil || resp.Events == nil {
		return addrs
	}

	for i := range resp.Events {
		if v := resp.Events[i].Kv.Value; v != nil {
			nodeData := ServiceInfo{}
			err := json.Unmarshal(v, &nodeData)
			if err != nil {
				logrus.Info("Parse node data error:", err)
				continue
			}
			addrs = append(addrs, nodeData)
		}
	}
	return addrs
}
