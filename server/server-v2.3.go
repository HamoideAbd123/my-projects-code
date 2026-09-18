package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/argosopentech/argos-translate-go"
)

func main() {
	// اول مرة بس فكي التعليق عشان ينزل موديل الترجمة
	// fmt.Println("بجلب موديل الترجمة ar<->en...")
	// argostranslate.DownloadAndInstall("ar", "en")
	// argostranslate.DownloadAndInstall("en", "ar")
	// fmt.Println("تم تحميل موديل الترجمة")

	// 1. يخدم ملفات html من نفس المجلد
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

	// 4. API الترجمة الاوفلاين
	http.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		text := r.URL.Query().Get("text")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if text == "" {
			http.Error(w, "نص فاضي", 400)
			return
	}

	translated := argostranslate.Translate(text, from, to)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"translated": translated})
	})

	fmt.Println("✅ السيرفر شغال على: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
