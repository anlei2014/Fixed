package main

import(
	"fmt"
	// "jedi-sim/jediSim"
	// "jedi-sim/msgHandler"
	// "jedi-sim/zmq4"
	"log"
	"net/http"
	"os"
	"sync"
	// "time"

	"jedi-sim/internal/handler"
	"jedi-sim/internal/model"
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
	wg.Add(1) // Adding two goroutines to the WaitGroup

	// Start ZMQ server goroutine
	//go zmqServer(&wg)

	// start event Machine to handle the request
	//go JediRun(&wg)

	// start HTTP server
	go startHTTPServer(&wg)
	// Wait for all goroutines to complete
	wg.Wait()
}

// func JediRun(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	var readHandler msgHandler.MSGHandler
// 	readHandler = jediSim.GetReadCanMSGHandler()
// 	sem := make(chan int)
// 	messageQueue := make(chan *[]int, 200)
// 	//isStart := make(chan int)
// 	//go sendInitMessage(isStart)

// 	go func() {
// 		//flag := true
// 		for {
// 			msg := readHandler.ReadMessage()
// 			//if flag {
// 			//	flag = false
// 			//	isStart <- 1
// 			//}
// 			for ix := 0; ix < len(msg); ix += 10 {
// 				m := msg[ix : ix+10]
// 				if m[0] != 0 {
// 					messageQueue <- &m
// 				}
// 			}
// 		}
// 		sem <- 1
// 	}()

// 	go func() {
// 		for {
// 			readHandler.UponMessage(*<-messageQueue)
// 		}
// 	}()

// 	<-sem

// }

// func zmqServer(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	// Initialize ZMQ context and socket
// 	context, _ := zmq4.NewContext()
// 	socket, _ := context.NewSocket(zmq4.REP)
// 	defer socket.Close()
// 	socket.Bind("ipc:///tmp/jedi-sim.sock") // IPC socket binding

// 	for {
// 		log.Println("ZMQ Server is is waiting for request from the client")

// 		// Receive message from ZMQ client
// 		msg, _ := socket.Recv(0)

// 		log.Println("Received ZMQ request:", msg)

// 		jediSim.ProcessZMQRequest(msg)

// 		// Send response back to ZMQ client
// 		socket.Send("Response from server", 0)
// 	}
// }

// 新增：启动 HTTP 服务
func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	// 加载错误码
	if err := model.LoadErrorCodes(); err != nil {
		log.Fatalf("Failed to load error codes: %v", err)
	}

	// 提供静态资源（注意路径配置）
	// pwd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal("无法获取当前目录:", err)
	// }

	// baseDir := fmt.Sprintf("%s/simulator/jedi", pwd)

	// log.Println("当前基础目录:", baseDir)
	// log.Println("静态资源目录:", fmt.Sprintf("%s/static", baseDir))

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
// func sendInitMessage(flag chan int) {
// 	<-flag
// 	var writeHandler = jediSim.GetWriteCanMSGHandler()

// 	/*
// 		writeHandler.SendMessage([]int{460, 1, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{527, 4, 6, 2, 11, 2, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{525, 16, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{525, 4, 6, 2, 10, 2, 0, 0, 0, 0})

// 		for i := 1; i < 15; i++ {
// 			writeHandler.SendMessage([]int{750, 3, i, 0, 0, 0, 0, 0, 0, 0})
// 		}
// 		writeHandler.SendMessage([]int{750, 3, 15, 81, 10, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{750, 3, 15, 209, 10, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{750, 3, 15, 209, 3, 0, 0, 0, 0, 0})

// 		for i := 2; i < 26; i++ {
// 			writeHandler.SendMessage([]int{1120, 1, i, 0, 0, 0, 0, 0, 0, 0})
// 		}

// 		writeHandler.SendMessage([]int{1134, 8, 3, 6, 1, 0, 0, 0, 0, 0})
// 		for i := 1; i < 19; i++ {
// 			writeHandler.SendMessage([]int{1130, 8, i, 0, 0, 0, 0, 0, 0, 0})
// 		}

// 		writeHandler.SendMessage([]int{1134, 8, 4, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{1134, 8, 3, 7, 0, 0, 0, 0, 0, 0})
// 	*/
// 	for i := 1; i < 29; i++ {
// 		writeHandler.SendMessage([]int{1130, 8, i, 0, 0, 0, 0, 0, 0, 0})
// 		time.Sleep(200 * time.Millisecond)
// 	}
// 	/*
// 		writeHandler.SendMessage([]int{1134, 8, 4, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{1134, 8, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{1134, 8, 7, 20, 19, 1, 10, 15, 10, 45})
// 		writeHandler.SendMessage([]int{520, 8, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{1134, 8, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{520, 8, 2, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{720, 0, 0, 0, 0, 0, 0, 0, 0, 0})
// 		writeHandler.SendMessage([]int{480, 7, 50, 80, 64, 6, 212, 48, 0, 0})
// 	*/
// }