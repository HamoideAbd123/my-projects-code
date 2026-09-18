package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	argostranslate "github.com/argosopentech/argos-translate-go"
)

type OllamaResponse struct {
	Response string `json:"response"`
}

const ollamaURL = "http://localhost:11434/api/generate"

// دالة الترجمة
func translate(text string, from string, to string) string {
	installedLanguages := argostranslate.GetInstalledLanguages()
	var fromLang, toLang argostranslate.Language
	
	for _, lang := range installedLanguages {
		if lang.Code == from { fromLang = lang }
		if lang.Code == to { toLang = lang }
	}
	if fromLang.Code == "" || toLang.Code == "" { return text }

	translation, err := argostranslate.Translate(text, fromLang, toLang)
	if err!= nil { return text }
	return translation
}

// دالة سؤال Ollama
func askOllama(prompt string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
	"model": "tinydolphin",
	"prompt": prompt,
		"stream": false,
	})
	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err!= nil { return "خطا في الاتصال بـ Ollama" }
	defer resp.Body.Close()

	var result OllamaResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Response
}

func main() {
	// نزل ملفات اللغة اول مرة
	fmt.Println("بتحمل ملفات الترجمة...")
	argostranslate.InstallLanguage("ar", "en")
	argostranslate.InstallLanguage("en", "ar")

	http.Handle("/", http.FileServer(http.Dir(".")))
	
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		qAr := r.URL.Query().Get("q")
		if qAr == "" {
			http.Error(w, "ضيف?q=السؤال", http.StatusBadRequest)
			return
	}

		// 1. عربي -> انجليزي
		qEn := translate(qAr, "ar", "en")
		
	// 2. نرسل للـ tinydolphin بالانجليزي + نخوفو عشان يرد مختصر
		prompt := fmt.Sprintf(`Answer this in English only, short and direct: %s`, qEn)
		ansEn := askOllama(prompt)

		// 3. انجليزي -> عربي
		ansAr := translate(ansEn, "en", "ar")
		ansAr = strings.TrimSpace(ansAr)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, ansAr)
	})

	fmt.Println("السيرفر شغال على http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
