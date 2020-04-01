#!/bin/bash

# set -x
OS=$(uname)
export GODEBUG=cgocheck=0

function build(){
    time go build -v
}
function install(){
    time go install -v
}

function fmt(){
    ls $1 | while read file; do
        file=$1/$file
        if [ -f $file ]; then
            go fmt $file
        else
            fmt $file
        fi
    done
}
# function go_publish(){
#     # ./warehouse/tools/bin/bindata/go-bindata.linux -o=vendor/github.com/mkevac/debugcharts/bindata/bindata.go -pkg=bindata vendor/github.com/mkevac/debugcharts/static/...
    
#     log_success "publish test asset"
#     # go get -u github.com/jteeuwen/go-bindata/...
#     # GOOS=linux GOARCH=amd64 go build -o go-bindata.linux
#     # GOOS=darwin GOARCH=amd64 go build -o go-bindata.darwin
#     # GOOS=windows GOARCH=amd64 go build -o go-bindata.exe
#     texe="go-bindata.exe"
#     if [ $OS = Linux ];then
#         texe="go-bindata.linux"
#     elif [ $OS = Darwin ];then
#         texe="go-bindata.darwin"
#     fi
#     ls ./warehouse/tools/bin/bindata/
#     fullexe="./warehouse/tools/bin/bindata/$texe"
    
#     log_success "bind test data"
#     [[ -f $fullexe ]] && $fullexe -o=$1 -pkg=$2 $3
# }