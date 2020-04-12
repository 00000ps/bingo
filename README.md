[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/00000ps/bingo) 

# bingo

bingo is a framework for software testing


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
# url = username:password@https://github.com/00000ps/bingo.git
```