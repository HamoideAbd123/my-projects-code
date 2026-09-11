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
	// 1. يخدم ملفات html
	http.Handle("/", http.FileServer(http.Dir(".")))

	// 2. API يجيب الموديلات من Ollama
	http.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://127.0.0.1:11434/api/tags")
		if err!= nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		bufio.NewReader(resp.Body).WriteTo(w)
	})

	// 3. API حق الشات مع stream
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		m := r.URL.Query().Get("model")
		
		modelName := "tinydolphin:latest"
		if m == "ar" { modelName = "tinydolphin-ar-v1:latest" }
		if m == "coder" { modelName = "dolphin-coder-v1:latest" }
		if m!= "" && m!= "ar" && m!= "coder" { modelName = m }

		requestBody, _ := json.Marshal(map[string]interface{}{
			"model": modelName,
			"prompt": q,
			"stream": true, // مهم: بنخليه true
		})

		resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
		if err!= nil {
			http.Error(w, "خطا في الاتصال بـ Ollama", 500)
			return
		}
		defer resp.Body.Close()

		// مهم جدا للـ stream
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Access-Control-Allow-Origin", "*")

	// نقرأ من Ollama سطر سطر ونرسلو للمتصفح طوالي
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			w.Write(line)
			w.Write([]byte("\n"))
			if f, ok := w.(http.Flusher); ok { // بنجبر السيرفر يرسل فورا
				f.Flush()
			}
	}
	})

	fmt.Println("السيرفر شغال على: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

