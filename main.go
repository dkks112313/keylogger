package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func writeToFile(st string) {
	f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(st)
	if err != nil {
		log.Fatal(err)
	}
}

func still(w http.ResponseWriter, r *http.Request) {
	save := make([]byte, 1024)
	n, err := r.Body.Read(save)
	if err != io.EOF && err != nil {
		log.Fatal(err)
	}
	st := string(save[:n])
	log.Print(st)

	if st == "favicon.ico" {
		return
	}
	writeToFile(st)
}

func main() {
	http.HandleFunc("/", still)
	http.ListenAndServe(":25583", nil)
}
