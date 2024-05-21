package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Title   string
	Content string
}

// テンプレートのキャッシュ
var templates = template.Must(template.ParseFiles("template.html"))

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

func main() {
	// ログファイルの設定
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)

	// 静的ファイルのサービング
	staticDir := "static"
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ルートハンドラの設定
	http.HandleFunc("/", handler)

	// サーバーの起動
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
