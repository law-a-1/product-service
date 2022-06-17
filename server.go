package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/law-a-1/product-service/ent"
	"github.com/law-a-1/product-service/ent/product"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	db     *ent.Client
	logger *zap.SugaredLogger
}

func NewServer(logger *zap.SugaredLogger, db *ent.Client) *Server {
	return &Server{
		router: chi.NewRouter(),
		db:     db,
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

//type productRequest struct {
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	Price       int    `json:"price"`
//	Stock       int    `json:"stock"`
//	Image       string `json:"image,omitempty"`
//	Video       string `json:"video,omitempty"`
//}

type userResponse struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	Role     string `json:"role"`
}

type errorResponse struct {
	Message string `json:"message"`
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
				JSON(w, http.StatusInternalServerError, nil, "failed to get all products")
				return
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

			JSON(w, http.StatusOK, productsResponse, "All Products fetched")
		})

		r.With(IsAuthorized, IsAdmin).Post("/", func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(5120)

			price, err := strconv.Atoi(r.FormValue("price"))
			if err != nil {
				panic("invalid price")
			}

			stock, err := strconv.Atoi(r.FormValue("stock"))
			if err != nil {
				panic("invalid stock")
			}

			// image, imageHeader, err := r.FormFile("image")
			// if err != nil {
			// 	log.Fatal(err)
			// 	panic("error image")
			// }
			// defer image.Close()

			// video, videoHeader, err := r.FormFile("video")
			// if err != nil {
			// 	panic("error video")
			// }
			// defer video.Close()

			// Create the uploads folder if it doesn't
			// already exist
			// err = os.MkdirAll("./uploads", os.ModePerm)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// // Create a new file in the uploads directory
			// imageDst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(imageHeader.Filename)))
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// defer imageDst.Close()

			// // Copy the uploaded file to the filesystem
			// // at the specified destination
			// _, err = io.Copy(imageDst, image)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// videoDst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(videoHeader.Filename)))
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// defer videoDst.Close()

			// // Copy the uploaded file to the filesystem
			// // at the specified destination
			// _, err = io.Copy(videoDst, video)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			_, err = s.db.Product.
				Create().
				SetName(r.FormValue("name")).
				SetDescription(r.FormValue("description")).
				SetPrice(price).
				SetStock(stock).
				SetImage("").
				SetVideo("").
				Save(r.Context())
			if err != nil {
				JSON(w, http.StatusInternalServerError, nil, "failed to create product")
				return
			}

			JSON(w, http.StatusCreated, nil, "Product created")
		})

		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					idString := chi.URLParam(r, "id")
					if idString == "" {
						JSON(w, http.StatusBadRequest, nil, "id cannot be empty")
						return
					}

					id, err := strconv.Atoi(idString)
					if err != nil {
						JSON(w, http.StatusBadRequest, nil, "invalid product id")
						return
					}

					p, err := s.db.Product.
						Query().
						Where(product.ID(id)).
						Only(r.Context())
					if err != nil {
						JSON(w, http.StatusNotFound, nil, "product not found")
						return
					}

					ctx := context.WithValue(r.Context(), "product", p)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			})

			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				p, ok := r.Context().Value("product").(*ent.Product)
				if !ok {
					JSON(w, http.StatusBadRequest, nil, "failed to parse product")
					return
				}

				pr := productResponse{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
					Stock:       p.Stock,
					Image:       p.Image,
					Video:       p.Video,
				}

				JSON(w, http.StatusOK, pr, "Product fetched")
			})

			r.Group(func(r chi.Router) {
				r.Use(IsAuthorized, IsAdmin)

				r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(5 << 20)

					price, err := strconv.Atoi(r.FormValue("price"))
					if err != nil {
						panic("invalid price")
					}

					stock, err := strconv.Atoi(r.FormValue("stock"))
					if err != nil {
						panic("invalid stock")
					}

					p, ok := r.Context().Value("product").(*ent.Product)
					if !ok {
						JSON(w, http.StatusInternalServerError, nil, "failed to parse product")
						return
					}

					_, err = p.Update().
						SetName(r.FormValue("name")).
						SetDescription(r.FormValue("description")).
						SetPrice(price).
						SetStock(stock).
						SetImage("").
						SetVideo("").
						SetUpdatedAt(time.Now()).
						Save(r.Context())
					if err != nil {
						JSON(w, http.StatusInternalServerError, nil, "failed to update product")
						return
					}
					JSON(w, http.StatusNoContent, nil, "Product updated")
				})

				r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
					p, ok := r.Context().Value("product").(*ent.Product)
					if !ok {
						JSON(w, http.StatusInternalServerError, nil, "failed to parse product")
						return
					}

					err := s.db.Product.
						DeleteOne(p).
						Exec(r.Context())
					if err != nil {
						JSON(w, http.StatusBadRequest, nil, "failed to delete product")
						return
					}

					JSON(w, http.StatusNoContent, nil, "Product deleted")
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

func JSON(w http.ResponseWriter, status int, v any, message string) error {
	w.WriteHeader(status)
	if status < 300 {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return err
		}
	} else {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}

	logType := "INFO"
	if status >= 300 {
		logType = "ERROR"
	}
	marshall, _ := json.Marshal(map[string]string{
		"type":    logType,
		"service": "products",
		"message": strconv.Itoa(status) + " - " + message,
	})
	req, _ := http.NewRequest("POST", os.Getenv("LOG_SERVICE_URL"), bytes.NewReader(marshall))
	req.Header.Add("Content-Type", "application/json")
	_, _ = http.DefaultClient.Do(req)

	return nil
}
