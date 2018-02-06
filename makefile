all:  mykv

test: *.go kv/kv.pb.go 
	go test -cover -coverprofile=c.out   ./ 
	go tool cover -html=c.out -o coverage.html

mykv:*.go  kv/kv.pb.go 
	go build -o mykv ./

kv/kv.pb.go: kv/kv.proto
	protoc -I kv/ kv/kv.proto --go_out=plugins=grpc:kv
