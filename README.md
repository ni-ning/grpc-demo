# gRPC Demo

## 代码生成

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/*.proto
```

## 服务端启动
```
go run server/server.go -p 8000
```

### 客户端启动
```
go run client/client.go -p 8000
```