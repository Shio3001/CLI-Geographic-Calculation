//get requestを受信SQL文を受け取る
//SQL Likeな文でgeojsonのデータを整理できるようにする
//url routing は

package main

import (
	handler "CLI-Geographic-Calculation/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Vercel と同じパスで叩けるようにする
	mux.HandleFunc("/api/", handler.Handler)

	addr := ":8080"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
