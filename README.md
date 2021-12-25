tgblock
===

命令行工具, 将telegram当成一个无限容量存储使用。

**NOTE: 无目录管理功能**

# 部署

## 服务端
### 使用docker-compose进行部署

```yml
version: "3.0"
services:
  tgblock:
    image: xxxsen/tgblock:v0.0.1
    restart: unless-stopped
    ports:
      - 127.0.0.1:8444:8444
    environment:
      - LISTEN=:8444   #服务监听地址
      - TOKEN=abc      #机器人token
      - MAX_FILE_SIZE=2147483648 #单个文件上传上限
      - BLOCK_SIZE=20971520 #单个文件块大小, 最大20M, 不可更改
      - CHATID=1234        #文件转发chatid, 上传文件后, 机器人会将文件转发给相应的chatid
      - LOG_LEVEL=trace
      - SECRETID=hello   # secretid, 相当于用户名
      - SECRETKEY=world  # secretkey, 相当于密码
      - DOMAIN=example.com  #对外暴露的域名, 用于拼接分享地址
      - SCHEMA=http         #对外暴露的schema, 用于拼接分享地址
      - CACHE_MEM_KEY_SIZE=10000  #缓存key数量(内存), 用于加速meta信息读取
      - CACHE_FILE_KEY_SIZE=500000 #缓存key数量(文件), 用于加速meta信息读取
      - TEMP_DIR=/tmp    #使用的临时目录
```

### 使用docker进行部署

根据需要自行替换对应的环境变量, 具体有哪些环境变量, 可以参考上面的docker-compose配置

```shell
docker run --rm -it -e "TOKEN=abc" -e "CHATID=12345" xxxsen/tgblock:v0.0.1
```

### 配置nginx

```conf
# 替换example为自己的域名
# 替换proxy_pass为自己的监听地址

server {
        listen 443 ssl;
        server_name example.com;
        proxy_connect_timeout 120s;
        proxy_read_timeout    120s;
        client_body_buffer_size  30M;

        location / {
               proxy_pass http://127.0.0.1:8444;
               proxy_http_version 1.1;
               proxy_set_header Host $http_host;
               proxy_set_header X-Real-IP $remote_addr;
               proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
               proxy_set_header X-Forwarded-Proto $scheme;
        }
} 
```

## 客户端

客户端是一个命令行工具, 自行从release下载对应系统/架构的二进制文件即可。

默认情况下, 需要手动创建客户端配置文件

windows配置文件位置为`c:\\tgblock\\client.json`

linux配置文件位置则是`/etc/tgblock/client.json`

## 相关的配置

```json
{
    "server":"https://example.com",
    "secret_id":"hello",
    "secret_key":"world",
    "max_sig_alive_time":60, 
    "max_file_size":2147483648,
    "block_size":20971520
}
```
|字段|类型|说明|
|---|---|---|
|server|string|服务端地址|
|secret_id|string|secret_id, 与服务端保持一致|
|secret_key|string|secret_key, 与服务端保持一致|
|max_sig_alive_time|int|sig有效时间, 默认即可, 防重放|
|max_file_size|int|最大文件大小, 需要与服务端一致, 填0则从服务端加载|
|block_size|int|最大文件块大小, 需要与服务端一致, 填0则从服务端加载|

# 使用

支持下面几个命令

|命令|说明|
|---|---|
|cat|读取一个文件|
|info|查看文件元信息|
|share|分享文件|
|upload|上传文件|
|download|下载文件|

具体的命令使用可以通过help进行查看
```shell
tgblock_cli ${cmd} --help

# tgblock_cli cat --help
```

## 一些具体的使用case

### 上传文件
```shell
# 上传文件main.go
./tgblock_cli upload -file=./main.go
```
```text
user@workstation:~/work/tgblock/cmd/cli$ ./tgblock_cli upload -file=./main.go
2021/12/25 21:33:25 upload succ, fileid:ABjscgVi0GMohLWhzUAJAADotAAwZSnV6nHwmAWxyli9_TOASsdccYbNAADkxAAUgACAQB
```

### 下载文件
```shell
# 下载指定fileid的数据, 存储为abc.txt
./tgblock_cli download -fileid=ABjscgVi0GMohLWhzUAJAADotAAwZSnV6nHwmAWxyli9_TOASsdccYbNAADkxAAUgACAQB -target=./abc.txt
```

```text
user@workstation:~/work/tgblock/cmd/cli$ ./tgblock_cli download -fileid=ABjscgVi0GMohLWhzUAJAADotAAwZSnV6nHwmAWxyli9_TOASsdccYbNAADkxAAUgACAQB -target=./abc.txt
2021/12/25 21:34:45 read file info succ, hash:db7e6ab616b5d1aa07e0dbbf6140f234, size:1153, block count:1
```
