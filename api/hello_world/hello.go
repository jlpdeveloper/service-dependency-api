package hello_world

import (
	"net/http"
)

func HelloWorld(rw http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	response := "hello world"
	if name != "" {
		response = "hello " + name
	}
	_, err := rw.Write([]byte(response))
	if err != nil {
		return
	}
}
