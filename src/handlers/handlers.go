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
	Question string `json:"question"`
}

type PostAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
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

func PostAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var pa PostAnswer
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &pa)
	if err != nil {
		http.Error(w, "Error parsing requst body", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=postgres password=postgres host=db port=5432 dbname=scalerz sslmode=disable")
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
	}
	defer conn.Close(ctx)

	queries := gen.New(conn)
	answer, err := queries.InsertAnswer(ctx, gen.InsertAnswerParams{
		Question: pa.Question,
		Answer:   pa.Answer,
	})
	if err != nil {
		http.Error(w, "Error inserting into database", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(answer)
}
