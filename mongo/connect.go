package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// addr []string{"139.155.56.188:27017"}
func NewConnect(addr []string, timeout time.Duration) (*mongo.Client, error) {
	var (
		client    *mongo.Client
		opts      *options.ClientOptions
		err       error
		ctx       context.Context
		cancelFun context.CancelFunc
	)

	opts = &options.ClientOptions{
		AppName:        nil,
		Auth:           nil,
		ConnectTimeout: &timeout,
		Hosts:          addr,
	}

	ctx, cancelFun = context.WithTimeout(context.Background(), timeout)
	defer cancelFun()
	// 建立连接
	if client, err = mongo.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return client, nil
}
