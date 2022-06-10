package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Listening for HTTP connections...")
    err := http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Shutting down")
}
