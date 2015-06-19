package main

import (
	"net/http"
	"ratelimit/rate"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Success!"))
}
func main() {

	singleHosted := rate.NewRateHandler(http.HandlerFunc(myHandler))

	http.ListenAndServe(":2705", singleHosted)

}
