package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"jedi-sim/internal/handler"
)

var state int
var mu sync.Mutex

type Message struct {
	Header string `json:"header"`
	Data   []int  `json:"data"`
}

func main() {
	// Enable line numbers and timestamp in logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	var wg sync.WaitGroup
	wg.Add(1)

	// start HTTP server
	go startHTTPServer(&wg)
	// Wait for all goroutines to complete
	wg.Wait()
}

// startHTTPServer starts the HTTP server
func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("无法获取当前目录:", err)
	}
	staticDir := fmt.Sprintf("%s/static", pwd)
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API 路由
	http.HandleFunc("/api/report", handler.ReportHandler)

	// 根路径重定向到 index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	// 启动服务
	log.Println("HTTP server starting on :9103")
	if err := http.ListenAndServe(":9103", nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
