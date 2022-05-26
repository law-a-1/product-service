package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"io"
	"net/http"
	"strings"
)

func (s Server) SetupMiddlewares() {
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.AllowContentType("application/json", "application/octet-stream", "multipart/form-data"))
	s.router.Use(cors.Default().Handler)
	s.router.Use(middleware.RequestLogger(&SugaredRequestLogger{Logger: s.logger}))

	s.router.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))

	s.router.Use(middleware.Recoverer)
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		splitToken := strings.Split(reqToken, "Bearer ")
		if splitToken[0] != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if splitToken[1] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reqToken = splitToken[1]

		req, err := http.NewRequest("GET", "https://auth-law-a1.herokuapp.com/user", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+reqToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var user userResponse
		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value("user").(userResponse)
		// TODO check if user is admin
		if u.Role != "user" {
			w.WriteHeader(http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//logger.Info(r.Method, r.URL.Path, r.Proto)
		next.ServeHTTP(w, r)
	})
}
