### description

golang micro-service framework

golang实现的微服务框架

### 推荐微服务架构分层

service接入层 -- controller业务逻辑层 -- model数据访问层

### install

```shell
go get github.com/pefish/go-core 
```

### road map
    1、自动生成controller模版代码
    2、实现服务间http请求看起来向本地调用一样的体验（像grpc一样编译成各语言平台代码）
    3、实现服务间rpc请求看起来向本地调用一样的体验
    4、自动生成api前置处理器模板
    5、移除iris

### 包依赖关系
```shell
service -> api-strategy -> builder -> session
service -> builder -> session
```

