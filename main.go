package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"weixin"

	"github.com/eperfect/goini"
)

var listenPort string
var AppID string
var AppKey string

func main() {
	initConfig()
	startServe()
}

func initConfig() {
	goini.InitConfig("config.ini")
	listenPort = goini.GetValue("", "port")
	AppID = goini.GetValue("", "app_id")
	AppKey = goini.GetValue("", "app_key")
}

func startServe() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", weixin.GetUserInfo(r))
	})

	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", weixin.GetSign())
	})

	go http.ListenAndServe(":8181", nil)
	acceptUserInput()
}

func acceptUserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		log.Println("command:", command)
		if command == "stop" {
			break
		}
	}
}
