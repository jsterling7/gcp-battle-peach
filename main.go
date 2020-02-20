package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/jsterling7/gcp-battle-peach/controller"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.Root)
	mux.HandleFunc("/battle", controller.Battle)


	const port = "8080"
	fmt.Println("Starting http server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}



