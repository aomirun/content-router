package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aomirun/content-router/api"
)

// httpHandler 是一个处理HTTP请求的处理器
func httpHandler(w http.ResponseWriter, r *http.Request) {
	// 创建一个缓冲区
	buf := api.NewBuffer()

	// 将请求数据写入缓冲区
	data := []byte("Hello, this is an HTTP request: " + r.URL.Path)
	buf.Write(data)

	// 创建路由器
	router := api.NewRouter()

	// 注册路由
	router.Match("Hello", func(ctx api.Context) error {
		response := "Processed: " + string(ctx.Buffer().Get())
		fmt.Fprintf(w, "%s", response)
		return nil
	})

	// 使用路由器处理请求
	_, err := router.Route(context.Background(), buf)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// 注册HTTP处理函数
	http.HandleFunc("/", httpHandler)

	fmt.Println("Server starting on :8080...")

	// 启动HTTP服务器
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
