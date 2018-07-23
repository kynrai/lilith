package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/handlers"
	"github.com/kynrai/lilith/internal/config"
	"github.com/kynrai/lilith/pkg/datastore_example"
	datastore_exampleH "github.com/kynrai/lilith/pkg/datastore_example/handlers"
	hello_worldH "github.com/kynrai/lilith/pkg/hello_world/handlers"
)

type Server struct {
	Router    *chi.Mux
	Datastore datastore_example.Repo
}

func New() *Server {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Use this config variable in some initialisation
	fmt.Println(conf.Secret)
	s := &Server{}

	s.Datastore = datastore_example.New()

	s.Router = chi.NewRouter()
	s.Router.Use(
		middleware.RedirectSlashes,
		middleware.DefaultCompress,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
		middleware.SetHeader("Content-Type", "application/json"),
	)

	s.Router.Method(http.MethodGet, "/health", Health())

	s.Router.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodGet, "/hello", hello_worldH.Hello())
		r.Method(http.MethodGet, "/hello/{name}", hello_worldH.HelloName())
		r.Method(http.MethodGet, "/datastore/{id}", datastore_exampleH.GetThing(s.Datastore))
		r.Method(http.MethodPost, "/datastore", datastore_exampleH.PutThing(s.Datastore))
	})

	return s
}

// Run the server
func (s *Server) Run() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	allowCredentials := handlers.AllowCredentials()
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT"})
	allowedOrigins := handlers.AllowedOriginValidator(allowedOriginValidator())
	corsHandler := handlers.CORS(allowedHeaders, allowCredentials, allowedMethods, allowedOrigins)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CompressHandler(corsHandler((s.Router)))))
}

func allowedOriginValidator() func(string) bool {
	r := regexp.MustCompile(`.*`)
	return func(s string) bool {
		return r.MatchString(s)
	}
}
