package main

import (
	"bingo/internal/app/server"
	"bingo/pkg/monitor"
)

func main(){
	server.Serve()
	monitor.Perform()
}