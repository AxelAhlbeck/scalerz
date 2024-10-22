package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scalerz/src/db/gen"

	"github.com/jackc/pgx/v5"
)

type PostQuestion struct {
	Sender   string `json:"sender"`
	Question string `json:"question"`
}

func PostQuestionHandler(w http.ResponseWriter, r *http.Request) {
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
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=postgres password=postgres host=db port=5432 dbname=scalerz sslmode=disable")
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
	}
	defer conn.Close(ctx)

	queries := gen.New(conn)
	answer, err := queries.GetAnswer(ctx, pq.Question)
	if err != nil {
		http.Error(w, "Error fetching from database", http.StatusInternalServerError)
	}
	if len(answer.Answer) > 0 {
		fmt.Fprintln(w, answer.Answer)
	} else {
		fmt.Println(w, "No answer to this question")
	}
}
