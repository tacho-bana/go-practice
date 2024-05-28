package main

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// ページデータの構造体
type PageData struct {
	Title       string
	Content     string
	Error       string
	SessionData map[string]interface{}
}

// テンプレートのキャッシュ
var templates *template.Template

func init() {
	// セッションに保存するために、map[string]interface{}のエンコードを登録
	gob.Register(map[string]interface{}{})

	// テンプレートのパスを指定してパース
	templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
}

// セッションデータをmap[string]interface{}に変換するヘルパー関数
func convertSessionValues(values map[interface{}]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for key, value := range values {
		if strKey, ok := key.(string); ok {
			converted[strKey] = value
		}
	}
	return converted
}

// メインページのハンドラ
func handler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	data := PageData{
		Title:       "Hello, Go!",
		Content:     "Welcome to the Go web server.",
		SessionData: convertSessionValues(session.Values),
	}
	err := templates.ExecuteTemplate(w, "template.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// フォームページのハンドラ
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "session-name")
		session.Values["name"] = r.FormValue("name")
		session.Values["message"] = r.FormValue("message")
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := templates.ExecuteTemplate(w, "form.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// エラーハンドラ
func errorHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Error",
		Content: "Something went wrong!",
	}
	err := templates.ExecuteTemplate(w, "error.html", data)
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
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Logging to file started")

	// 静的ファイルのサービング
	staticDir := filepath.Join("static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ルートハンドラの設定
	http.HandleFunc("/", handler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/error", errorHandler)

	// ログミドルウェアを使用
	loggedRouter := loggingMiddleware(http.DefaultServeMux)

	// サーバーの起動
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
