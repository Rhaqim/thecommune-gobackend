package main

import (
	"github.com/Rhaqim/thecommune-gobackend/pkg/handlers"
)

func main() {
	// r := handlers.NewRouter()
	// fmt.Println("Server is running on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", r))

	run := handlers.GinRouter()

	run.Run(":8088")

}
