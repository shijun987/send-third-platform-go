# send-third-platform-henan

## 编译

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

$env:GOOS="linux"
$env:GOARCH="amd64"
go build
```

## 构建镜像

```shell
docker build -t registry.cn-hangzhou.aliyuncs.com/whxph/send-third-platform-henan .

docker push registry.cn-hangzhou.aliyuncs.com/whxph/send-third-platform-henan
```
