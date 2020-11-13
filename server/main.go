package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.iu.edu/evogelsa/strain-sense/routes"
)

const (
	ROOT_PATH = "/wearables"
	STATIC    = ROOT_PATH + "/static"
	PORT      = ":32321"
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

// newRouter creates a new mux router with the desired paths and structures
func newRouter() *mux.Router {
	// create new router
	router := mux.NewRouter()

	// create static file system for the style sheet and other static assets
	fs := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	router.PathPrefix(STATIC).Handler(http.StripPrefix(STATIC, fs))

	// define all routes
	router.HandleFunc(ROOT_PATH+"/login", routes.DisplayLogin).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/login", routes.AuthenticateLogin).Methods("POST")
	router.HandleFunc(ROOT_PATH+"/create", routes.DisplayCreate).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/create", routes.CreateAccount).Methods("POST")
	router.HandleFunc(ROOT_PATH+"/dashboard", routes.DisplayDashboard).Methods("GET")

	return router
}

func main() {
	router := newRouter()

	log.Fatal(http.ListenAndServe(
		PORT,
		router,
	))
	// log.Fatal(http.ListenAndServeTLS(
	//     PORT,
	//     "/etc/letsencrypt/live/ethanvogelsang.xyz/cert.pem",
	//     "/etc/letsencrypt/live/ethanvogelsang.xyz/privkey.pem",
	//     router,
	// ))
}
