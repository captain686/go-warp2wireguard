# 🤠`go-warp2wireguard`

> 提取 `warp`+ 节点信息生成`wireguard`或`clash`节点配置

- 自动注册`warp`+ 账号
- 自动刷取`warp`+ 流量
- 自动节点测速

## 💿使用

```
./go-warp2wireguard -h
Usage:
  -t string
        operating mode [ wireguard | clash ] (default "wireguard")
```

> `clash subscribe url` http://ip:8888

### [`warp 代理类型查看`](https://www.cloudflare.com/cdn-cgi/trace)