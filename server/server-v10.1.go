package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. يخدم كل ملفات html من المجلد الحالي
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// 2. API يجيب الموديلات من Ollama
	http.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		resp, err := http.Get("http://127.0.0.1:11434/api/tags")
		if err!= nil {
			http.Error(w, "Ollama ما شغال: "+err.Error(), 500)
			return
	}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		bufio.NewReader(resp.Body).WriteTo(w)
	})

	// 3. API حق الشات مع stream
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		q := r.URL.Query().Get("q")
		m := r.URL.Query().Get("model")
		if m == "" { m = "tinydolphin:latest" }

		requestBody, _ := json.Marshal(map[string]interface{}{
			"model": m,
			"prompt": q,
			"stream": true,
		})

		resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err!= nil {
			http.Error(w, "خطا في الاتصال بـ Ollama. شغلو: ollama serve", 500)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			w.Write(line)
			w.Write([]byte("\n"))
			if f, ok := w.(http.Flusher); ok { f.Flush() }
		}
	})

	fmt.Println("✅ السيرفر شغال على: http://127.0.0.1:8080")
	fmt.Println("✅ افتحي: http://127.0.0.1:8080/home-v10.1.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
