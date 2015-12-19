package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eperfect/goini"
)

func main() {
	startServer()
}

func startServer() {
	initServerConfig()
	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", GetSign())
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", GetUserInfo(r))
	})
	log.Fatal(http.ListenAndServe(":8181", nil))
}

func initServerConfig() {
	goini.InitConfig("config.ini")
	AppID = goini.GetValue("", "app_id")
	AppKey = goini.GetValue("", "app_key")
	fmt.Println(AppID)
	fmt.Println(AppKey)

}
