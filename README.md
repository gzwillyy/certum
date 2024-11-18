


## Certum 验证服务器

用于验证 Certum 证书的轻量级 HTTP 服务器。此服务器创建一个验证文件，并根据证书验证的需要通过特定的 URL 路径提供该文件。

- 在 `.well-known/pki-validation/certum.txt` 处创建一个验证文件。
- 在 HTTP 上提供验证文件。
- 在收到 `Ctrl+C` (SIGINT) 时处理正常关机。
- 终止时自动清理临时文件。

```
curl -s https://raw.githubusercontent.com/gzwillyy/certum/master/install.sh | bash && chmod +x ./certum_validation && ./certum_validation

```

## 优化服务器TCP负载均衡性能

适用于在TCP负载均衡的节点上执行 以提高性能

```
curl -sSL https://github.com/gzwillyy/certum/raw/master/optimize_system.sh | bash
```

## 配置宝塔编译nginx， 配置IP证书 及 基于IP的 TCP + TLS 访问

适用于在TCP负载均衡的节点上执行 以配置访问方式

```
curl -L -o build_tcp_conf https://github.com/gzwillyy/certum/raw/master/build_tcp_conf && chmod +x ./build_tcp_conf && ./build_tcp_conf
```


## bt 安装脚本


- 通用
```sh
if [ -f /usr/bin/curl ];then curl -sSO http://io.bt.sb/install/install_panel.sh;else wget -O install_panel.sh http://io.bt.sb/install/install_panel.sh;fi;bash install_panel.sh
```

- debian ubuntu
```
wget -O install.sh http://io.bt.sb/install/install-ubuntu_6.0.sh && bash install.sh
```

- centos
```
yum install -y wget && wget -O install.sh http://io.bt.sb/install/install_6.0.sh && sh install.sh
```

- 升级为企业
```
curl https://io.bt.sb/install/update_panel.sh|bash
```