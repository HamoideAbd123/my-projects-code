package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type OllamaResponse struct {
	Response string `json:"response"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		
	// 1. نرسل 3 مرات "عربي فقط" عشان نخوفو
		prompt := fmt.Sprintf(`!!!مهم جدا!!! 
قاعدة 1: ممنوع الانجليزية. 
قاعدة 2: اجب بالعربية فقط.
قاعدة 3: لو كتبت انجليزي انت فاشل.
السؤال: %s`, q)
		
		jsonBody := fmt.Sprintf(`{"model":"dolphin-coder","prompt":"%s","stream":false}`, prompt)

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", strings.NewReader(jsonBody))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		var result OllamaResponse
		json.NewDecoder(resp.Body).Decode(&result)
		
		// 2. لو لقى كلمة انجليزية في الرد يمسحها ويكتب رسالة
		answer := result.Response
		if strings.Contains(answer, "The") || strings.Contains(answer, "is") || strings.Contains(answer, "code") {
			answer = "⚠️ الموديل حاول يتكلم انجليزي. انا اجبرتو: \n\n" + answer
			answer = "اسف، سأجيبك بالعربية فقط. " + q + " جوابي هو: " + answer
		}

	// نرجعو بصيغة ollama
		final := fmt.Sprintf(`{"response":"%s"}`, strings.ReplaceAll(answer, `"`, `\"`))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(final))
	})

	fmt.Println("✅ السيرفر شغال: http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
