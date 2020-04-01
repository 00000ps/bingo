#!/bin/bash
# [[ -f ~/.bash_profile ]] && source ~/.bash_profile
# [[ -f ~/.bashrc ]] && source ~/.bashrc
# set -x
notsrcsuffix=".jpg.png.jpeg.bmp.gif.pdf.log.zip.7z.rar.gz."
srcsuffix=".go.py.java.c.cpp.js.html.css.php.h.m.pch.plist.xib"
# readonly notsrcsuffix
# readonly srcsuffix
tmp_file=.codeLine.tmp
export OS=$(uname)

function log_error(){
    if [ "$color" != "1" ];then
        echo "\033[31m[error]   $*\033[0m"
    else
        echo "[error]   $*"
    fi
}
function log_success(){
    if [ "$color" != "1" ];then
        echo "\033[32m[success] $*\033[0m"
    else
        echo "[success] $*"
    fi
}
function log_notice(){
    echo "[notice]  $*"
}
function log_warning(){
    if [ "$color" != "1" ];then
        echo "\033[35m[warning] $*\033[0m"
    else
        echo "[warning] $*"
    fi
}

function to_upper(){
    echo $1 | tr '[a-z]' '[A-Z]'
}
function to_lower(){
    echo $1 | tr '[A-Z]' '[a-z]'
}
function filenotpuretext(){
    echo $notsrcsuffix | grep -c ".$(to_lower $1)"
}
function filesuffix(){
    echo ${1##*.}
}
function trimspace(){
    # echo $1
    echo $1 | sed 's/^[ \t]*//g' |sed 's/[ \t]*$//g'
}


function list(){
    ls $1 | while read line
    do
        full=$1/$line
        test -f $full
        if [ $? -eq 0 ]; then
            #echo "$2 $full"
            wc -l $full
        else
            #list "$full" "$2"
            list "$full"
        fi
    done
}


function code_check(){
    codeLine=$2
    [[ A$2 = A ]] && codeLine=$(cat $tmp_file)
    #log_warning $codeLine $1 $2

    ls $1 | while read file; do
        file=$1/$file
        #file=$(pwd)/$1/$file
        #log_notice $file
        #wc -l $file|awk '{print $1}'
        #test $file
        if [ -f $file ]; then
            sf=$(filesuffix $file)
            if [ $(filenotpuretext $sf) -ne 0 ];then
                log_warning "$file is not source code file"
            else
                line=`wc -l $file|awk '{print $1}'`
                #codeLine=$line
                [[ -e $tmp_file ]] && codeLine=$[$(cat $tmp_file)+$line]
                echo $codeLine > $tmp_file
                log_notice "file: $file; line: $line; total: $codeLine"
            fi
        else
            code_check $file
            log_warning "dir: $file"
        fi
    done
}
function code(){
    code_check $1 0
    log_success "code line: $(cat $tmp_file)"
    [[ -f $tmp_file ]] && rm $tmp_file
}

function md5_sum(){
    if [ "$OS" = "Linux" ];then
        md5sum $1|awk '{print $1}'
    elif [ "$OS" = "Darwin" ];then
        #alias md5sum="md5"
        md5 $1|awk '{print $4}'
    fi
}
function md5_check(){
    if [ -f $1 ]; then
        echo "$(md5_sum $1)\t$1"
    else
        ls $1 | while read file; do
            file=$1/$file
            if [ -f $file ]; then
                echo "$(md5_sum $file)\t$file"
            else
                md5_check $file
            fi
        done
    fi
}
function MD5(){
    md5_file=md5.log
    md5_check $1|sort > $md5_file
    awk '{print $1}' $md5_file | uniq -d | sort -rn | while read line; do
        grep $line $md5_file | awk -v var=$line '{print $2}END{print "---- "var"  "NR" ----\n"}'
    done
    #log_success "code line: $(cat $tmp_file)"
    exit
}

function func_usage(){
    echo -e "reset               # sync source code with remote by force"
    echo -e "code path           # calculate code line"
    echo -e "MD5 path            # check duplicated file"
}