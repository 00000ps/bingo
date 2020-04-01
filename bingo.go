package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	help = flag.Bool("h", false, "help")
	debug = flag.Bool("debug", false, "debug mode")
	verbose = flag.Bool("v", false, "verbose")
)

func main(){
	flag.Parse()
	// flag.Usage = flag.Usage()

	if *help{
		flag.Usage()
		os.Exit(1)
	}
	if *debug {
		// TODO: make debug mode.
	}
	if *verbose {
		// TODO: set log to output verbose message.
	}

	fmt.Println("bingo")
	fmt.Println("this is a project of testing framework that design for server testing")
	
	e := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	e.GET("/x", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
		"messgae":"success",
	})})
	e.Run()
}