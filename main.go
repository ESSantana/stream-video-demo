package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Handle("/videos/*", http.StripPrefix("/videos", http.FileServer(http.Dir("./videos/output/63205490-0165-429b-a517-776ae0bde6ae/"))))
	router.Handle("/page", http.StripPrefix("/page", http.FileServer(http.Dir("./static"))))

	defer func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			panic(err)
		}
	}()
	
	fmt.Println("Server is running on port 8080")
}
