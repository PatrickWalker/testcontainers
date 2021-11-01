package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	fmt.Println("Server started and listening")
	http.ListenAndServe(":8080", nil)
}

//HelloServer is a helloworld function i stole from the internet. I feel remorse.
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I'm getting requests")
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
