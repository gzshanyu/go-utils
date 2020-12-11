module github.com/gzshanyu/go-utils

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/shomali11/util v0.0.0-20200329021417-91c54758c87b
	github.com/shopspring/decimal v1.2.0
	github.com/valyala/fasthttp v1.18.0 // indirect
	go.mongodb.org/mongo-driver v1.3.5
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/grpc v1.30.0 // indirect
)

// 替换为v1.26.0版本的gRPC库
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
