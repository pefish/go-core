### description

golang web framework

### hierarchy

service access layer -- business logic layer -- data access layer

### roadmap

    1、generate automaticly controller template code
    2、abstract http and rpc request
    3、generate automaticly api strategy template code
    4、remove iris

### dependencies relationship

```shell
service -> api-strategy -> builder -> session
service -> builder -> session
```

