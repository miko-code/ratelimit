package main

import (
	"net/http"
	"ratelimit/rate"
)

func myHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("DO SOME API STUFF!"))
}

func main() {
	r := http.NewServeMux()
	r.Handle("/bla", rate.RateHandler(http.HandlerFunc(myHandler)))
	http.ListenAndServe(":2705", r)

}
