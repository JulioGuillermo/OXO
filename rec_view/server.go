package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile("index.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = w.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
	})*/
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile("../record.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = w.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
	})

	err := http.ListenAndServe("0.0.0.0:9393", nil)
	if err != nil {
		fmt.Println(err)
	}
}
