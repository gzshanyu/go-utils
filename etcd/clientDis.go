package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
)

type ClientDis struct {
	connect     *Connect
	serviceList map[string]string
	lock        sync.Mutex
}

// addr: []string{"8.210.85.242:2379"}
func NewClientDis(addr []string) (*ClientDis, error) {
	var (
		err       error
		clientDis *ClientDis
		connect   *Connect
	)

	if connect, err = NewConnect(addr, 5000); err != nil {
		return nil, err
	}

	clientDis = &ClientDis{
		connect:     connect,
		serviceList: make(map[string]string),
	}

	return clientDis, nil
}

func (c *ClientDis) GetService(prefix string) ([]string, error) {
	var (
		addrs []string
		err   error
		resp  *clientv3.GetResponse
	)

	if resp, err = c.connect.client.Get(context.TODO(), prefix, clientv3.WithPrefix()); err != nil {
		return nil, err
	}

	addrs = c.ExtractAddrs(resp)
	go c.watcher(prefix)
	return addrs, nil
}

func (c *ClientDis) SetService(key, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.serviceList[key] = val
}

func (c *ClientDis) DelService(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.serviceList, key)
}

func (c *ClientDis) ExtractAddrs(resp *clientv3.GetResponse) []string {
	var addrs = make([]string, 0)
	if resp == nil || resp.Kvs == nil || len(resp.Kvs) == 0 {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			c.SetService(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}

func (c *ClientDis) watcher(prefix string) {
	var rch = c.connect.client.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				c.SetService(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				c.DelService(string(ev.Kv.Key))
			}
		}
	}
}

func (c *ClientDis) SerList2Array() []string {
	c.lock.Lock()
	defer c.lock.Unlock()
	var addrs = make([]string, 0)

	for _, v := range c.serviceList {
		addrs = append(addrs, v)
	}

	return addrs
}
