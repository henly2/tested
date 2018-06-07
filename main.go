package main

import (
	"time"
	"net/http"
	"os"
	"flag"
	"github.com/gin-gonic/gin"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Data struct{
	Msg string `json:"msg"`
}

var cmd = flag.String("c", "old", "mode")

func main()  {
	if *cmd == "old" {
		go startOld()
	} else if *cmd == "new" {
		go startNew()
	}
}

func test(w http.ResponseWriter, req *http.Request)  {
	// get argv
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	w.Header().Set("Connection", "close")
	w.Write(b)
}

func startOld()  {
	http.Handle("/test", http.HandlerFunc(test))

	server := &http.Server{
		Addr:        ":8023",
		Handler:     nil,
		ReadTimeout: time.Second*20,
		IdleTimeout:time.Second*5,
	}

	err := server.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}
}

func startNew() {
	engine := gin.Default()
	engine.POST("/test", func(ctx *gin.Context) {
		d := Data{}
		ctx.ShouldBindJSON(&d)
		if d.Msg == "ok" {
			ctx.String(http.StatusOK, "ok")
		} else {
			ctx.String(http.StatusOK, "err")
		}
	})

	engine.Run(":8023")
}

func startClient(count int)  {
	d := Data{
		Msg:"ok",
	}

	dd,_ := json.Marshal(d)
	ddd, err := httpPost("http://127.0.0.1:8023/test", dd)
	fmt.Println(string(ddd), err)
}

func httpPost(path string, data []byte) ([]byte, error) {
	client := http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		ResponseHeaderTimeout: time.Second * 30,
	}}

	resp, err := client.Post(path, "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}