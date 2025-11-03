package helloworld

import (
	"io"
	"log"
	"net/http"
)

func HelloWorld(rw http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	response := "hello world"
	if name != "" {
		response = "hello " + name
	}
	_, err := io.WriteString(rw, response)
	if err != nil {
		log.Println(err)
	}
}
