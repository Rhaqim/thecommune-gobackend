package main

import (
	"fmt"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/handlers"
	_ "github.com/lib/pq"
)

func main() {
	r := handlers.NewRouter()
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)

}
