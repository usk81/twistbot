package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

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
	r.Post("/bot", func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("can not get request"))
		}
		defer r.Body.Close()
		bs, _ := ioutil.ReadAll(r.Body)
		w.Write(bs)
	})

	logRoutes(r, logger)
	return r
}
