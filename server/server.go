package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kynrai/lilith/config"
	"github.com/kynrai/lilith/services/datastore_example"
	"github.com/kynrai/lilith/services/hello_world"
)

type Server struct {
	Router    *mux.Router
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

	s.Router = mux.NewRouter()
	s.Router.StrictSlash(true)

	s.Router.Handle("/health", Health()).Methods(http.MethodGet)

	v1 := s.Router.PathPrefix("/v1").Subrouter()
	v1.Handle("/hello", hello_world.Hello()).Methods(http.MethodGet)
	v1.Handle("/hello/{name}", hello_world.HelloName()).Methods(http.MethodGet)
	v1.Handle("/datastore/{id}", datastore_example.GetThing(s.Datastore)).Methods(http.MethodGet)
	v1.Handle("/datastore", datastore_example.PutThing(s.Datastore)).Methods(http.MethodPost)
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
