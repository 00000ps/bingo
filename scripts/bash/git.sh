#!/bin/bash

export Branch=$(git branch|grep '*'|awk '{print $2}')

function branch(){
    echo "$Branch"
}
function reset(){
    git fetch --all && git reset --hard origin/${Branch}
}
function git_init(){
    git config --global user.email $2              ## 改为你的邮箱
    git config --global user.name $1               ## 改为你的账号
    git config --global remote.origin.push 'refs/heads/*:refs/for/*'
}