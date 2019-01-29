package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//var movies []Movie

func main() {

	// movies = append(movies, Movie{Name: "Batman Begins", Year: 2013, Director: "Alguien"})
	// movies = append(movies, Movie{Name: "Rápido y furioso 8", Year: 2017, Director: "de los mejores"})
	// movies = append(movies, Movie{Name: "El señor de los anillos", Year: 2017, Director: "Ni puta idea"})

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index).Methods("GET")

	router.HandleFunc("/peliculas", MovieList).Methods("GET")
	router.HandleFunc("/peliculas/{id}", MovieShow).Methods("GET")

	router.HandleFunc("/peliculas", MovieAdd).Methods("POST")
	router.HandleFunc("/peliculas/{id}", MovieUpdate).Methods("PUT")
	router.HandleFunc("/peliculas/{id}", MovieDelete).Methods("DELETE")

	server := http.ListenAndServe(":8080", router)
	panic(server)
}
