module young

go 1.14

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.8 // indirect
	github.com/jonboulle/clockwork v0.2.1 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/genproto v0.0.0-20200914193844-75d14daec038 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20200915000551-32e7099cccac // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
