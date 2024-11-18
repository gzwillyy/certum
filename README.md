


## Certum 验证服务器

```
curl -s https://raw.githubusercontent.com/gzwillyy/certum/master/install.sh | bash
```

## 优化服务器TCP负载均衡性能

```
curl -sSL https://github.com/gzwillyy/certum/raw/master/optimize_system.sh | bash
```

## 配置宝塔编译nginx， 配置IP证书 及 基于IP的 TCP + TLS 访问

```
curl -L -o build_tcp_conf  https://github.com/gzwillyy/certum/raw/master/build_tcp_conf

chmod +x ./build_tcp_conf

./build_tcp_conf
```