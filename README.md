# fire-hizard-detect/fire-hizard-detect



## 框架
Go Zero

### 生成带缓存的mysql代码
```shell
goctl model mysql ddl -src="./sql/*.sql" -dir="./model" -c --style goZero
```
### 基于api生成代码
```shell
goctl api go --api api.api --dir ./ --style goZero
```
### 基于proto生成代码
```shell
goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. --style goZero
```

## 模块
* `service` 各服务
* `common` 公共模块

### `service`下模块
* `user` 用户模块+权限认证(小程序, 微信等等)

## 中间件
* mysql
* redis
* etcd

## 备注
先以http实现基本登录, 后续引入grpc和etcd拆分模块