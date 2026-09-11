package main
import "net/http"
import "log"

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	log.Println("السيرفر شغال على http://127.0.0.1:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
