package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	port := ":3000"
	if envPort, ok := os.LookupEnv("PORT"); ok {
		port = fmt.Sprintf(":%s", envPort)
	}

	srv := &http.Server{
		Addr:         port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	staticFs := RedirectingFileSystem{http.Dir("static")}
	http.Handle("/", http.FileServer(staticFs))

	log.Fatal(srv.ListenAndServe())
}
