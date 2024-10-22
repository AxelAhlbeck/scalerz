package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"scalerz/src/handlers"
	mw "scalerz/src/handlers/middlewares"
)

type PostQuestion struct {
	Sender   string `json:"sender"`
	Question string `json:"question"`
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlers.PostQuestionHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

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