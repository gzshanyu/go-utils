package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type ServiceReg struct {
	connect       *Connect
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFunc    context.CancelFunc
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

// addr: []string{"8.210.85.242:2379"}
func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	var (
		serReg  *ServiceReg
		err     error
		connect *Connect
	)

	if connect, err = NewConnect(addr, timeNum*1000); err != nil {
		return nil, err
	}

	serReg = &ServiceReg{
		connect: connect,
	}

	if err = serReg.SetLease(timeNum); err != nil {
		return nil, err
	}

	go serReg.listenLeaseRespChan()
	return serReg, nil
}

// 设置租约
func (s *ServiceReg) SetLease(timeNum int64) error {
	var (
		ctx           context.Context
		cancelFunc    context.CancelFunc
		err           error
		lease         clientv3.Lease
		leaseResp     *clientv3.LeaseGrantResponse
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	)

	lease = clientv3.NewLease(s.connect.client)

	// 设置租约时间
	if leaseResp, err = lease.Grant(context.TODO(), timeNum); err != nil {
		return err
	}

	// 设置续租
	ctx, cancelFunc = context.WithCancel(context.TODO())
	if keepAliveChan, err = lease.KeepAlive(ctx, leaseResp.ID); err != nil {
		return nil
	}

	s.cancelFunc = cancelFunc
	s.lease = lease
	s.leaseResp = leaseResp
	s.keepAliveChan = keepAliveChan
	return nil
}

// 监听续租情况
func (s *ServiceReg) listenLeaseRespChan() {
	for {
		select {
		case keepResp := <-s.keepAliveChan:
			if keepResp == nil {
				fmt.Printf("租约已失效")
				goto END
			} else {
				// 每秒会续租一次, 所以就会受到一次应答
				fmt.Println("收到自动续租应答:", keepResp.ID)
			}
		}
	}
END:
}

// 通过租约注册服务
func (s *ServiceReg) PutService(key, val string) error {
	var (
		err error
		kv  clientv3.KV
	)
	kv = clientv3.NewKV(s.connect.client)
	if _, err = kv.Put(context.TODO(), key, val, clientv3.WithLease(s.leaseResp.ID)); err != nil {
		return err
	}

	return nil
}

// 撤销租约
func (s *ServiceReg) RevokeLease() error {
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	_, err := s.lease.Revoke(context.TODO(), s.leaseResp.ID)
	return err
}
