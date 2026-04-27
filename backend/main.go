package main

import (
	"log"
	"net/http"
	"stonkx/database"
	h "stonkx/handler"
)

func main() {
	pool, err := database.Connect("dbconfig.env")
	if err != nil {
		log.Fatal("connecting db: %w", err)
	}

	hand := h.New(pool)

	http.HandleFunc("/api/getcharts", hand.GetData)
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	//http.HandleFunc("/", homeHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
