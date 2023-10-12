### icp-cli
使用golang 实现的备案信息查询命令行工具，使用cobra框架实现。

⚠️⚠️⚠️⚠️⚠️    

不要开多线程并发查询，会被封ip的。想要解决这个，自己去看源码改一下就好了。

⚠️⚠️⚠️⚠️⚠️

### 依赖
需要opencv支持,https://github.com/hybridgroup/gocv#macos 自行参考
本人在mac 上编译，通过brew 安装如下依赖,其他平台还没测试，回头再继续补充文档

```bash
brew install pkgconfig
brew install opencv
```

### 安装
目前还没有发布二进制文件，可以自行编译安装
```bash
cd icp-cli && make buildLocal
```

### 使用
看帮助命令就好了
```bash
icp-cli -h
```

### 示例
```bash
icp-cli check -u baidu.com
```

### TODO
- [x] 实现 verbose 的flag,增加详细输出
- [x] 查询接口输出到指定的文件，也就是 oyaml,ojson,ocsv 三个flag
- [x] 目前通过公司名称查询的时，只能查询到40条数据，后续需要循环查处所有的数据
- [x] 增加API方式调用
- [x] 嵌入前端GUI页面
- [x] 增加docker 方式部署


### 免责声明

本软件仅用于学习和研究使用,用户须自行承担使用本软件引起的所有法律和相关责任。

本软件不提供任何形式的担保,不保证特定用途下的适用性和安全性。使用本软件可能面临网络安全、信息安全等风险,用户须自行承担所有风险。

本声明的解释、效力、执行等相关事宜,均适用中华人民共和国法律的规定。因使用本软件引起的纠纷,均应由用户自行负责处理,与软件作者无关。