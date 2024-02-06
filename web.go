package main

import "net/http"

func initweb() {
	http.HandleFunc("/data/2.5/weather", handlerRoot)
	http.ListenAndServe(listenHost, nil)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
	// aaa := r.FormValue("aaa")
	// if aaa != "" {
	// 	w.Write([]byte(aaa))
	// }
}
