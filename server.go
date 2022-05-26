package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/law-a-1/product-service/ent"
	"github.com/law-a-1/product-service/ent/product"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	router *chi.Mux
	db     *ent.Client
}

func NewServer(db *ent.Client) *Server {
	return &Server{
		router: chi.NewRouter(),
		db:     db,
	}
}

func (s Server) SetupMiddlewares() {
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8"))
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.Recoverer)
}

type productsResponse struct {
	Products []productResponse `json:"products"`
	Count    int               `json:"count"`
}

type productResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
	Video       string `json:"video"`
}

type productRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image,omitempty"`
	Video       string `json:"video,omitempty"`
}

type userResponse struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	Role     string `json:"role"`
}

func (s Server) SetupRoutes() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	s.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	s.router.Route("/products", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			products, err := s.db.Product.
				Query().
				All(r.Context())
			if err != nil {
				panic(err)
			}

			var productsResponse productsResponse
			for _, p := range products {
				productResponse := productResponse{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
					Stock:       p.Stock,
					Image:       p.Image,
					Video:       p.Video,
				}
				productsResponse.Products = append(productsResponse.Products, productResponse)
			}
			productsResponse.Count = len(productsResponse.Products)

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(productsResponse); err != nil {
				panic(err)
			}
		})

		r.With(isAuthorized, isAdmin).Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req productRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			_, err := s.db.Product.
				Create().
				SetName(req.Name).
				SetDescription(req.Description).
				SetPrice(req.Price).
				SetStock(req.Stock).
				SetImage(req.Image).
				SetVideo(req.Video).
				Save(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusCreated)
		})

		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					idString := chi.URLParam(r, "id")
					if idString == "" {
						w.WriteHeader(http.StatusBadRequest)
					}

					id, err := strconv.Atoi(idString)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
					}

					p, err := s.db.Product.
						Query().
						Where(product.ID(id)).
						Only(r.Context())
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					ctx := context.WithValue(r.Context(), "product", p)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			})

			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				p, ok := r.Context().Value("product").(*ent.Product)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(productResponse{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
					Stock:       p.Stock,
					Image:       p.Image,
					Video:       p.Video,
				}); err != nil {
					panic(err)
				}
			})

			r.With(isAuthorized, isAdmin).Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
				var req productRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					panic(err)
				}

				p, ok := r.Context().Value("product").(*ent.Product)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
				}

				w.WriteHeader(http.StatusNoContent)
				_, err := p.Update().
					SetName(req.Name).
					SetDescription(req.Description).
					SetPrice(req.Price).
					SetStock(req.Stock).
					SetImage(req.Image).
					SetVideo(req.Video).
					SetUpdatedAt(time.Now()).
					Save(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			})

			r.With(isAuthorized, isAdmin).Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
				p, ok := r.Context().Value("product").(*ent.Product)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
				}

				err := s.db.Product.
					DeleteOne(p).
					Exec(r.Context())
				if err != nil {
					panic(err)
				}
				w.WriteHeader(http.StatusNoContent)
			})
		})
	})
}

func isAuthorized(next http.Handler) http.Handler {
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

		defer res.Body.Close()

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

func isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value("user").(userResponse)
		// TODO check if user is admin
		if u.Role != "user" {
			w.WriteHeader(http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}

func (s Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		return err
	}
	return nil
}
