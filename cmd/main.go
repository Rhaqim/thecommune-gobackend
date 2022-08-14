package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rhaqim/thecommune-gobackend/pkg/handlers"
)

func main() {
	r := handlers.NewRouter()
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
