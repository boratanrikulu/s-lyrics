package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/boratanrikulu/s-lyrics/controllers"
	"github.com/boratanrikulu/s-lyrics/controllers/api"
)

func main() {
	godotenv.Load()

	err := cloneOrPullDatabase(os.Getenv("DATABASE_ADDRESS"))
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.WelcomeGet).Methods("GET")
	r.HandleFunc("/logout", controllers.LogoutGet).Methods("GET")
	r.HandleFunc("/spotify", controllers.SpotifyGet).Methods("GET")
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/search", api.Search).Methods("POST")
	// api.HandleFunc("/write", api.WriteLyrics).Methods("POST")

	serve(r, "3000")
}

func serve(r *mux.Router, defaultPort string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort // Default port if not specified
	}

	log.Println("Server started at :" + port + ".")
	log.Fatalln(http.ListenAndServe(":"+port, r))
}

func cloneOrPullDatabase(databaseURL string) error {
	var errOutput bytes.Buffer

	if folderExists("./database/.git") {
		cmd := exec.Command("git", "pull", "origin", "master")
		cmd.Dir = "./database"
		cmd.Stderr = &errOutput

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("%s", errOutput)
		}

		return nil
	}

	cmd := exec.Command("git", "clone", databaseURL, "./database")
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s", errOutput)
	}

	return nil
}

func folderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
