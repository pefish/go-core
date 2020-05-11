
#### 描述

本实例演示了

1、如何构建一个微服务

2、如何查询mysql

3、微服务分层架构的最佳实践

4、如何自动生成swagger文档

#### 启动

```shell
GO_CONFIG=`pwd`/config/local.yaml GO_SECRET=`pwd`/secret/local.yaml go run ./bin/test/
```

#### 测试

```shell
curl localhost:3000/api/test/v1/test_api?user_id=1
```

#### 生成swagger

```shell
go run scripts/gene_swagger.go
```
