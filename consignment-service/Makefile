build:
# 一定要注意 Makefile 中的缩进，否则 make build 可能报错 Nothing to be done for build
# protoc 命令前边是一个 Tab，不是四个或八个空格
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/shippy/consignment-service proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build 
	docker build -t consignment-service .
run:
	# 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
	# 可添加 -d 参数将微服务放到后台运行
	docker run -p 50051:50051 consignment-service