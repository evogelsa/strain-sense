package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

	"github.iu.edu/evogelsa/strain-sense/routes"
)

const (
	ROOT   = "/wearables"
	STATIC = ROOT + "/static"
	PORT   = ":32321"
)

// neuteredFileSystem is a file system that prevents directory listing
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open checks if the requested file is a directory or not
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

// initCache creates a connections to the redis server to store cookies
func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}

	routes.Cache = conn
}

// newRouter creates a new mux router with the desired paths and structures
func newRouter() *mux.Router {
	// create new router
	router := mux.NewRouter()

	// create static file system for the style sheet and other static assets
	fs := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	router.PathPrefix(STATIC).Handler(http.StripPrefix(STATIC, fs))

	// define all routes
	router.HandleFunc(ROOT+"/login", routes.DisplayLogin).Methods("GET")
	router.HandleFunc(ROOT+"/login", routes.AuthenticateLogin).Methods("POST")
	router.HandleFunc(ROOT+"/create", routes.DisplayCreate).Methods("GET")
	router.HandleFunc(ROOT+"/create", routes.CreateAccount).Methods("POST")
	router.HandleFunc(ROOT+"/dashboard", routes.DisplayDashboard).Methods("GET")
	router.HandleFunc(ROOT+"/dashboard", routes.SendUserData).Methods("POST")
	router.HandleFunc(ROOT+"/dashboard/log", routes.LBPLog).Methods("POST")

	return router
}

func main() {
	// initialize the cache connection
	initCache()
	// load in the credential files or create one if non existant
	routes.InitCredentials()

	// make local router for testing
	router := newRouter()
	log.Fatal(http.ListenAndServe(
		PORT,
		router,
	))

	// make https router for server
	// log.Fatal(http.ListenAndServeTLS(
	//     PORT,
	//     "/etc/letsencrypt/live/ethanvogelsang.xyz/cert.pem",
	//     "/etc/letsencrypt/live/ethanvogelsang.xyz/privkey.pem",
	//     router,
	// ))
}
