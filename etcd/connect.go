package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

type Connect struct {
	config clientv3.Config
	client *clientv3.Client
}

// addr etcd集群地址: []string{"8.210.85.242:2379"}
// timeout 连接超时时间: 单位毫秒
func NewConnect(addr []string, timeout int64) (*Connect, error) {
	var (
		err  error
		conf clientv3.Config
		conn *Connect
		cli  *clientv3.Client
	)
	// 配置客户端
	conf = clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Duration(timeout) * time.Millisecond,
	}
	// 建立连接
	if cli, err = clientv3.New(conf); err != nil {
		return nil, err
	}

	conn = &Connect{
		config: conf,
		client: cli,
	}

	return conn, nil
}
