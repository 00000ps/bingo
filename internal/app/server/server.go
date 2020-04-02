package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Serve() {
	go server1(8080)
	go server2(8081)

	select{}
}

func server1(p int){
	e := gin.New()
	gin.SetMode(gin.DebugMode)

	// g := e.Group("/api/v1")
	// g.Any("/*action", proxyServ)
	e.Any("/",proxyServ)

	err := e.Run(":"+strconv.Itoa(p))
	if err != nil {
		fmt.Println(err)
	}
}
func proxyServ(c *gin.Context){
		host := c.Query("host")
		// host = "www.google.com"
		fmt.Println("redirected to", host)
		c.Request.Header.Add("requester-uid", "id")
		p := httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = host
				req.Host = host
			},
		}
		p.ServeHTTP(c.Writer, c.Request)
}

func server2(p int) {
	//被代理的服务器host和port
	h := &handle{host: "127.0.0.1", port: "80"}
	err := http.ListenAndServe(":"+strconv.Itoa(p), h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
type handle struct {
	host string
	port string
}
func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse("http://" + h.host + ":" + h.port)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
