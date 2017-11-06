package main

import (
	"github.com/yJepo/Asssignment2_cloud/ass2"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/exchange/evaluationtrigger", resource.FullTriggerCheck) //GET
	http.HandleFunc("/exchange", resource.HandlerPost)                        //POST
	http.HandleFunc("/exchange/", resource.HandlerGetDel)                     //GET and DELETE
	http.HandleFunc("/exchange/latest", resource.HandlerLatest)               //POST (and GET)
	http.HandleFunc("/exchange/average", resource.HandlerAverage)             //POST
	http.ListenAndServe(":"+port, nil)
}
