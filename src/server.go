package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	mw "scalerz/middlewares"
	"strings"
)

type PostQuestion struct {
	Sender   string `json:"sender"`
	Question string `json:"question"`
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePostQuestion(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func handlePostQuestion(w http.ResponseWriter, r *http.Request) {
	var pq PostQuestion
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &pq)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}
	var response string
	if strings.HasSuffix(string(pq.Question), "?") {
		response = "good question"
	} else {
		response = "thanks for telling me that"
	}
	fmt.Fprintln(w, response)
}

func main() {
	portPtr := flag.String("port", "8081", "The port the server listens on.")
	flag.Parse()
	finalHandler := mw.LoggingMiddleware(
		mw.RecoveryMiddleware(
			mw.CORSMiddleware(
				mw.AuthenticationMiddleware(
					mw.RateLimitingMiddleware(
						mw.RequestIDMiddleware(questionHandler),
					),
				),
			),
		),
	)

	http.HandleFunc("/question", finalHandler)
	fmt.Println("Server is running at http://localhost:" + *portPtr)
	log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
}
