package main
import ("fmt"; "net/http"; "bytes"; "encoding/json"; "log")

type OllamaResponse struct{Response string `json:"response"`}

func ask(model, prompt string) string {
	body,_:=json.Marshal(map[string]interface{}{"model":model,"prompt":prompt,"stream":false})
	resp,_:=http.Post("http://localhost:11434/api/generate","application/json",bytes.NewBuffer(body))
	defer resp.Body.Close(); var r OllamaResponse; json.NewDecoder(resp.Body).Decode(&r); return r.Response
}

func main(){
	http.HandleFunc("/ask",func(w http.ResponseWriter,r *http.Request){
		q := r.URL.Query().Get("q")
		m := r.URL.Query().Get("model") // بنقرا الموديل من الرابط
		
	// بنحول الكلمة القصيرة لاسم الموديل الكامل
		modelName := "tinydolphin:latest" // الافتراضي انجليزي
		if m == "ar" { modelName = "tinydolphin-ar:latest" }
		if m == "coder" { modelName = "dolphin-coder:latest" }
		
		ans := ask(modelName, q)
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		fmt.Fprint(w, ans)
	})
	fmt.Println("السيرفر شغال: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}
