//http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func startWebserver() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	go http.ListenAndServe(":8080", nil)
}

func startLoadTest() {
	count := 0
	for {
		resp, err := http.Get("http://google.com/")
		// resp, err := http.Get("http://localhost:8080/")

		if err != nil {
			panic(fmt.Sprintf("Got error: %v", err))
		}
		//#1 solution
		//???server needs to close request to drain request body???
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		fmt.Println(resp)
		log.Printf("Finished GET request #%v", count)
		count++
	}

}

func main() {

	startWebserver()

	startLoadTest()
}
