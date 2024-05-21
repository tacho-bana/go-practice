package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PageData struct {
	Title   string
	Content string
}

// テンプレートのキャッシュ
var templates = template.Must(template.ParseFiles(filepath.Join("templates", "template.html")))

// ハンドラ関数
func handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Hello, Go!",
		Content: "Welcome to the Go web server.",
	}
	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ログミドルウェア
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %s %d %s\n",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			time.Since(start),
		)
	})
}

func main() {
	// ログファイルの設定
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)

	// 静的ファイルのサービング
	staticDir := filepath.Join("static")
	fs := http.FileServer(http.Dir(staticDir))

	// ルートハンドラの設定
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)

	// ログミドルウェアを使用
	loggedRouter := loggingMiddleware(http.DefaultServeMux)

	// サーバーの起動
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
