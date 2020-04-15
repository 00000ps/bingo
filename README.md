[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/00000ps/bingo) 

# bingo
bingo is a framework for software testing

## 代码结构

```bash
- [x] ├── LICENSE            ### MIT开源协议
- [x] ├── README.md          ### Readme
- [x] ├── _config.yml        ### 环境配置文件
- [x] ├── bingo.go           ### 主程序入口
+ [x] ├── cmd                ### 工具目录，每个子目录应是对应一个独立工具
+ [x] │   └── debug          ### 调试代码，用于debug
- [x] ├── conf               ### 配置文件目录[toml]
- [x] │   └── bingo.toml     ### 主配置文件
+ [x] ├── doc                ### 文档[md]
+ [ ] │   └── DEV.md         ### 开发计划
- [x] ├── go.mod             ### modules依赖库
- [x] ├── go.sum             ### modules依赖库及版本
- [ ] ├── god.sh             ### shell控制脚本，如编译、打包等
+ [x] ├── internal           ### 内部代码，禁止bingo外部代码引用
+ [x] │   ├── app            ### 主程序代码 
+ [x] │   │   ├── frame      ### 框架代码
+ [x] │   │   └── server     ### 内部server
+ [x] │   └── pkg            ### 主（内部）程序库
+ [x] │       └── testing    ### 内部测试库，定制化封装
- [x] ├── pkg                ### 公共库
- [x] │   ├── cache          ### 内存型缓存库
- [x] │   ├── cmd            ### 命令行库
- [x] │   ├── cv             ### 计算机视觉库
- [x] │   ├── datas          ### 数据结构库
- [ ] │   ├── encode         ### 编解码/序列化/反序列化库
- [ ] │   ├── format         ### 格式化库
- [ ] │   ├── log            ### 日志库
- [x] │   ├── net            ### 网络库
- [ ] │   ├── perf           ### 并发控制及调度库
- [ ] │   ├── ps             ### 进程及环境库
- [ ] │   ├── testing        ### 公共测试库
- [ ] │   └── utils          ### 基础库
- [x] └── scripts            ### 脚本目录
- [x]     └── bash           ### shell脚本[sh]
```


# 初始化
## go安装及配置
```bash
# 以下以linux示例，go建议安装1.14及以上版本，下载地址 https://studygolang.com/dl
    
# 安装golang
ver="1.14.1"
wget https://studygolang.com/dl/golang/go${ver}.linux-amd64.tar.gz
tar zxvf go${ver}.linux-amd64.tar.gz

# 设置GOPATH
export GOPATH=$(pwd)
export PATH=$PATH:$GOPATH/bin

# 设置go proxy，解决墙的问题
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

## git配置
```bash
git config --global user.name "John Doe"
git config --global user.email johndoe@example.com
```

## github
```bash
vim .git/config
# url = https://github.com/00000ps/bingo.git
# 修改为
# url = https://username:password@github.com/00000ps/bingo.git
```