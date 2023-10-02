package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"idbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getmovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
func deletemovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	found := false
	for idx, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			found = true
			break
		}
	}
	if found {
		w.WriteHeader(http.StatusOK)
		message := "Movie deleted successfully"
		response := map[string]interface{}{
			"message": message,
			"movies":  movies,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
		message := "Movie not found"
		response := map[string]string{"message": message}
		json.NewEncoder(w).Encode(response)
	}

}

func getmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	id := strconv.Itoa((len(movies) + 1))
	movie.ID = id
	movies = append(movies, movie)
	w.WriteHeader(http.StatusOK)
	message := "Movie Created Successfully"
	response := map[string]interface{}{
		"message": message,
		"movie":   movie,
	}
	json.NewEncoder(w).Encode(response)
}

func updatemovie(w http.ResponseWriter, r *http.Request) {
	// in this func we will actully delete the old movie and then add the new movie
	w.Header().Set("Contect-Type", "application/json")
	params := mux.Vars(r)
	for idx, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			id := params["id"]
			movie.ID = id
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{"1", "676767", "Movie One", &Director{"John", "Arora"}})
	movies = append(movies, Movie{"2", "696f97", "Movie Two", &Director{"Steve", "Arora"}})

	r.HandleFunc("/movies", getmovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getmovie).Methods("GET")
	r.HandleFunc("/movies", createmovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updatemovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deletemovie).Methods("DELETE")

	fmt.Println("Starting Server at port 8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
