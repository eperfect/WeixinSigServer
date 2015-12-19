package main

import (
	"net/http"
	"fmt"
	"log"

	"weixin"
	"github.com/eperfect/goini"
)

var listenPort string

func main() {
	initConfig()
	startServe()
}

func initConfig() {
	goini.InitConfig("config.ini")
	listenPort = goini.GetValue("", "port")
	weixin.AppID = goini.GetValue("", "app_id")
	weixin.AppKey = goini.GetValue("", "app_key")
}

func startServe() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", weixin.GetUserInfo(r))
	})

	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", weixin.GetSign(r))
	})
	
	log.Fatal(http.ListenAndServe(":" + listenPort, nil))
}
