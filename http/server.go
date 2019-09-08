package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/tcotav/proxyblock" pb
	"github.com/elazarl/goproxy"
)

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	MaxMinCount := flag.Int("max", 100, "Max count per minute allowed")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	hitCounter := pb.NewCountData(*MaxMinCount)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Print(formatLog(req))
		if hitCounter.ShouldBlock() {
			log.Print("blocked", formatLog(req))
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusTooManyRequests, fmt.Sprintf("Number of hits exceeded threshold of %d", *MaxMinCount))
		}
		return req, nil
	})
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

type CountServer struct {
	ListenPort  string
	MaxMinCount int
}

func formatLog(req *http.Request) string {
	return fmt.Sprintf("%s %s %s", req.RequestURI, req.RemoteAddr, req.Proto)
}


