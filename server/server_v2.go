package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OllamaResponse struct {
	Response string `json:"response"`
}

func ask(model, prompt string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"prompt": prompt,
		"stream": false,
	})
	resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err!= nil {
		return "خطا في الاتصال بـ Ollama"
	}
	defer resp.Body.Close()
	var r OllamaResponse
	json.NewDecoder(resp.Body).Decode(&r)
	return r.Response
}

func main() {
	// 1. ده عشان يخدم ملف ال html بتاعك
	http.Handle("/", http.FileServer(http.Dir(".")))

	// 2. ده API حق الشات
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		m := r.URL.Query().Get("model")

		// بنحدد الموديل حسب الاختيار
		modelName := "tinydolphin:latest"
		if m == "ar" {
			modelName = "tinydolphin-ar-v1:latest"
		}
		if m == "coder" {
			modelName = "dolphin-coder-v1:latest"
		}
		if m == "en" {
			modelName = "tinydolphin:latest"
		}

		ans := ask(modelName, q)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]string{"response": ans})
	})

	fmt.Println("السيرفر شغال على: http://127.0.0.1:8080/home.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
