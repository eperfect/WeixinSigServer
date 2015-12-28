package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eperfect/goini"
)

var (
	AppID      string
	AppKey     string
	listenPort string
)

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
		fmt.Fprintf(w, "%q", GetUserInfo(r))
	})

	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "%q", GetSign(r))
	})

	http.ListenAndServe(":"+listenPort, nil)
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
