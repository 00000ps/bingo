#!/bin/bash

# set -x 
# [[ -f ~/.bash_profile ]] && source ~/.bash_profile
# [[ -f ~/.bashrc ]] && source ~/.bashrc
CURRENT=bingo
readonly CURRENT

# [[ ! $var ]] && echo "var unset"
# [[ ! -n "$var" ]] && echo "var unset"
# test -z "$var" && echo "var unset"
# [[ "$var" = "" ]] && echo "var unset"
# source scripts/bash/*.sh
utilslibfile="./scripts/bash/utils.sh"
readonly utilslibfile
[[ -f $utilslibfile ]] && source $utilslibfile

golibfile="./scripts/bash/go.sh"
readonly golibfile
[[ -f $golibfile ]] && source $golibfile

gitlibfile="./scripts/bash/git.sh"
readonly gitlibfile
[[ -f $gitlibfile ]] && source $gitlibfile

step=0
fstart=`date +%s`


function heap(){
    wget http://xxx:8060/debug/pprof/heap
    go tool pprof heap # svg/png/help top/top10/list/web
    # go tool pprof -inuse_space http://xxx:8060/debug/pprof/heap
    # go tool pprof -alloc_space http://xxx:8060/debug/pprof/heap

    # go run -gcflags '-m -l' bingo.go # escape analyze
    # go build -gcflags '-m -l' main.go 或者 go build -gcflags ‘-m -m’ main.go    前者消除内联了，
    # go build -gcflags '-m -l' main.go -l 一个，表示消除内敛
    # go build -gcflags '-m -l -l' main.go -l 两个 ，表示内联级别比默认强
    # go build -gcflags '-m -l' main.go -l 3个，强内敛，二进制包体积变大，但是不稳定，可能有bug

    # /debug/pprof/profile：访问这个链接会自动进行 CPU profiling，持续 30s，并生成一个文件供下载
    # /debug/pprof/heap： Memory Profiling 的路径，访问这个链接会得到一个内存 Profiling 结果的文件
    # /debug/pprof/block：block Profiling 的路径
    # /debug/pprof/goroutines：运行的 goroutines 列表，以及调用关系

    # 使用 web funcName 的方式，只打印和某个函数相关的内容
    # 运行 go tool pprof 命令时加上 --nodefration=0.05 参数，表示如果调用的子函数使用的 CPU、memory 不超过 5%，就忽略它，不要显示在图片中

    # go tool trace trace.out
}
function publish(){
    # log_notice "begin to publish..." " testshop/data -> testshop/asset/asset.go"
    # go_publish testshop/asset/asset.go asset testshop/data/...
    
    log_notice "begin to building..." " command: time go build -v"
    build
    if [ $? -ne 0 ];then
        # cd ${GOPATH}/src && rm $CURRENT
        time go build -v
	    if [ $? -ne 0 ];then
            time go install -v 2>&1 | tee -a build.log
            sed -i -e ':a;N;s/\n/ /g;ta' build_log
            log_notice "build failed" " details: `cat build_log`"      
        fi
	    exit 1
    else
        # log_notice "building ${CURRENT} for linux" " details: GOOS=linux GOARCH=amd64 time go build -v -o ${CURRENT}.linux"
        # GOOS=linux GOARCH=amd64 time go build -v -o ${CURRENT}.linux
        # log_notice "building ${CURRENT} for darwin" " details: GOOS=darwin GOARCH=amd64 time go build -v -o ${CURRENT}.darwin"
        # GOOS=darwin GOARCH=amd64 time go build -v -o ${CURRENT}.darwin
        # log_notice "building ${CURRENT} for windows" " details: GOOS=windows GOARCH=amd64 time go build -v -o ${CURRENT}.exe"
        # GOOS=windows GOARCH=amd64 time go build -v -o ${CURRENT}.exe

        chmod +x bingo* supervise
        [[ -d output ]] && rm -rf output/*
        [[ -d output ]] || mkdir -p output
        cp -r README.md $CURRENT* conf warehouse doc lib supervise ./output
        # log_success "zipping to output..."
        # log_notice "build success, zipping to output..." " details: zip -r output.zip output"
        # zip -r output.zip output
        log_success "=====list all files===="
        log_notice "build success" " build finished.\n $(md5sum bingo)"
        #mkdir -p output/src && cp -r lib testshop tools *.go ./output/src
        
        # cd ${GOPATH}/src && rm $CURRENT
        # return 0
    fi
}
function cases(){
    git status| grep cases.go$ | awk '{print $2}' | while read line; do
        git show HEAD:$line > tmp
        diff tmp $line | grep '>' | awk -F'(' '{print $2}'| awk -F')' '{print $1}'| awk -F'_' '{print "go run bingo.go drun -d "$NF" -url \"http://127.0.0.1\" || echo "$NF" >> result"}' >> ids.sh
        rm tmp
    done
    cat ids.sh && sh ids.sh && rm ids.sh
    cat result && rm result
}

function log(){
    msg=`grep "_$1(t " -rn testshop/*|grep func|grep ' *super.Testing'`
    [[ "a$msg" == "a" ]] && msg=`grep "_$1_" -rn testshop/*|grep func|grep ' *super.Testing'`
    [[ "a$msg" == "a" ]] && msg=`grep " $1(t " -rn testshop/*|grep func|grep ' *super.Testing'`
    # echo "raw: $msg"
    raw=`echo $msg|awk '{print $2}'|awk -F'(' '{print $1}'`
    # echo "found case: $raw"
    [[ "a$msg" == "a" ]] && echo "case $1 not found" && exit 1
    msg=`echo $msg|awk -F':' '{print $1":"$2}'`
    echo "found case $raw in file: $msg"
    file=$(echo $msg|awk -F':' '{print $1}')
    line=$(echo $msg|awk -F':' '{print $2}')
    id=$(git blame $file|grep " $line) "|awk '{print $1}')
    # git show $id --stat| awk '{print $0}'
    # git show $id --stat
    git show $id --stat|head -5
}
function pb(){
    [ $# -lt 3 ] && exit 1
    file=$2
    module=$3
    export PATH=$(pwd)/warehouse/tools/bin/protobuf:$PATH
    chmod +x ./warehouse/tools/bin/protobuf/*
    pname=$(grep '^package ' $file |awk '{print $2}'|awk -F ';' '{print $1}')
    dir="./testshop/$module/$pname"
    [[ $module = $pname ]] && dir="./testshop/$module"
    [[ -d $dir ]] || mkdir -p $dir; cp $file $dir
    cd $dir
    protoc --go_out=. *.proto
    
    sed -i 's/github.com\/lib\/3rd\/github.com/g' *.pb.go
    cd -
}
function case_list(){
    grep 'NewCase' -wrn testshop/*|grep -w 'New' |awk -F')' '{print $1" "$3}'|awk -F'(' '{print $1" "$2" "$3}'|awk -F'/' '{print $1"/"$2"/"$3"/"$4" "$5}'|awk '{print $1""$4"*"$7}'|awk -F'*' '{print $1" "$2" "$3}' > log/caselist
}
function change(){
    log_success "prepare cases list"
    case_list
    log_success "change case type name"
    cat log/caselist | while read line; do
        pkg=`echo $line|awk '{print $1}'`
        name=`echo $line|awk '{print $2}'`
        feature=`echo $name|awk -F'_' '{print $1}'`
        id=`echo $line|awk '{print $3}'`
        
        #newname="${feature}_${id}"
        newname="${name}_${id}"
        
        log_success "changing: $name -> ${newname}"
        grep -rnw $name $pkg/* |awk -F':' '{print $1}'|sort|uniq|while read file; do
            grep ' NewCase()' $file|grep -q '.New(' && sed -i '/NewCase()/d' $file
            sed -i s"/\<${name}\>/${newname}/g" $file
            #sed -i s"type /${name} struct/type ${name}_${id} struct/g" $file
        done
    done
}
function duplicated(){
    case_list
    cat log/caselist |awk '{print $3}'|sort|uniq -c|sort -n
}
function eva(){
    name="${CURRENT}"
    # [[ -f $GOPATH/bin/$name ]] || name="${CURRENT}.exe"
    # echo "args: $*"
    # echo "cmds: go build -tags=opencv -v && ./$name eva $*"
    go build -tags=opencv -v && ./$name eva $*
}


function usage(){
    echo -e "update              # pull the latest source code"
    echo -e "pb file module      # auto generate proto interface source code"
    
    func_usage
    # echo -e "reset               # sync source code with remote by force"
    
}
function case_name(){
    file_name=$1
    cat $file_name |grep "super.Testing"  | grep "func" |awk '{print $2}' |awk -F[\(] '{print "testsuite.AddCase("$1")"}'
}
function case_name_by_file(){
    file_list_name=$1
    cat $file_list_name |while read line
    do
        case_name $line
    done
}
function case_seq(){
    file_name=$1
    cat $file_name |grep "super.Testing"  | grep "func" |awk '{print $2}' |awk -F[_\(] '{printf $3","}'
}
function case_seq_by_file(){
    file_list_name=$1
    cat $file_list_name |while read line
    do
        case_seq $line
    done
}
function rename_realID(){
    file_name=$1
}
function rename_realID_by_file(){
    file_list_name=$1
    cat $file_list_name |while read line
    do
        rename_realID $line
    done
}

if [ $# -eq 0 ];then
    # 无参数则进行编译， sh god.sh    ->    function 'build' (in god.sh) = go build
    build
elif [ "$1" = "style" ];then
    # sh god.sh style    ->    function 'style' (in god.sh)
    style
else
    # sh god.sh otherCmd    ->    function 'otherCmd' (in god.sh) 
    $* 2>/dev/null
    if [ $? -ne 0 ];then
        if [ -d "cmd/$1" ];then
            go run cmd/$1/*
        else
            # sh god.sh otherCmd    ->    go build && ./bingo otherCmd
            build $*
        fi
    fi
fi
