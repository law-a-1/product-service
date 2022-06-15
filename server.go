package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/law-a-1/product-service/ent"
	"github.com/law-a-1/product-service/ent/product"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	router *chi.Mux
	db     *ent.Client
	cache  *redis.Client
	logger *zap.SugaredLogger
}

func NewServer(logger *zap.SugaredLogger, db *ent.Client, rd *redis.Client) *Server {
	return &Server{
		router: chi.NewRouter(),
		db:     db,
		cache:  rd,
		logger: logger,
	}
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

type decrementStockRequest struct {
	Amount int `json:"amount"`
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

		r.With(IsAuthorized, IsAdmin).Post("/", func(w http.ResponseWriter, r *http.Request) {
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

			r.Group(func(r chi.Router) {
				r.Use(IsAuthorized, IsAdmin)

				r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
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

				r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
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

				r.Post("/{id}/decrement-stock", func(w http.ResponseWriter, r *http.Request) {
					p, ok := r.Context().Value("product").(*ent.Product)
					if !ok {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					var req decrementStockRequest
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					if p.Stock < req.Amount {
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					w.WriteHeader(http.StatusOK)
					_, err := p.Update().
						SetStock(p.Stock - req.Amount).
						SetUpdatedAt(time.Now()).
						Save(r.Context())
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				})
			})
		})
	})
}

func (s Server) Start() error {
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), s.router); err != nil {
		return err
	}
	return nil
}
