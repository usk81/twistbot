package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

// TwistOutgoingRequest is requests from Twist to outgoing webhook endpoint
type TwistOutgoingRequest struct {
	EventType         string `schema:"event_type"`
	WorkspaceID       int    `schema:"workspace_id"`
	Content           string `schema:"content"`
	UserID            int    `schema:"user_id"`
	UserName          string `schema:"user_name"`
	URLCallback       string `schema:"url_callback"`
	URLTTL            int    `schema:"url_ttl"`
	MessageID         int    `schema:"message_id"`
	ThreadID          int    `schema:"thread_id"`
	ThreadTitle       string `schema:"thread_title"`
	ChannelID         int    `schema:"channel_id"`
	ChannelName       string `schema:"channel_name"`
	CommentID         int    `schema:"comment_id"`
	ConversationID    int    `schema:"conversation_id"`
	ConversationTitle string `schema:"conversation_title"`
	VerifyToken       string `schema:"verify_token"`
}

var decoder = schema.NewDecoder()

func router(logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(zapLogger(logger))
	r.Use(middleware.Recoverer)

	// handlers
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Post("/bot", botHandler)
	r.Post("/check_request", checkRequestHandler)

	logRoutes(r, logger)
	return r
}

func botHandler(w http.ResponseWriter, r *http.Request) {
	lg, _ := zap.NewProduction()
	defer lg.Sync()

	if r == nil {
		lg.Error("can not get request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can not get request"))
	}
	if err := r.ParseForm(); err != nil {
		lg.Warn(fmt.Sprintf("failed to parse request %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request is invalid"))
	}
	var req TwistOutgoingRequest
	if err := decoder.Decode(&req, r.PostForm); err != nil {
		lg.Error(fmt.Sprintf("fail to decode parsed request %s : %#v", err.Error(), r.PostForm))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(req.Content))
}

func checkRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can not get request"))
	}
	defer r.Body.Close()
	bs, _ := ioutil.ReadAll(r.Body)
	w.Write([]byte(fmt.Sprintf("header: %s, body: %s", r.Header.Get("Content-Type"), string(bs))))
}
