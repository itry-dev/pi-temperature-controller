package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

var ctx = context.Background()
var validHTTPOrigins = regexp.MustCompile("(http|https)://(localhost:3000|127.0.0.1:3000)")
var useCors = false

type errResponse struct {
	Err  string
	Code uint64
}

func main() {
	http.HandleFunc("/api", handler)
	log.Fatal(http.ListenAndServe("localhost:5001", nil))
	log.Print("server is listening on port 5001")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Headers %v", r.Header)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, access-control-allow-origin")

	if useCors {
		origins := validHTTPOrigins.FindStringSubmatch(r.Header.Get("Origin"))
		if origins != nil {
			w.Header().Set("Access-Control-Allow-Origin", origins[1]+"://"+origins[2])
		}
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		log.Fatal(err)
	}
	w.Write(out)

}

func getJSONErrResponse(err error) []byte {
	jerr := errResponse{}
	jerr.Err = err.Error()

	jObject, _ := json.Marshal(err)

	return jObject
}

func doWriteErrHeader(w http.ResponseWriter, err error) {
	//log.Print(err.Error())
	w.WriteHeader(http.StatusBadRequest)
	w.Write(getJSONErrResponse(err))
}
