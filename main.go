package kvs

import (
	"fmt"
	"log"
	"net/http"
	//"github.com/yourusername/kvs" // replace with your actual module path
)

func main() {
	store := NewStore()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		store.HandleWebSocket(w, r)
	})

	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}