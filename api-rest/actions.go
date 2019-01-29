package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/*Acá tendremos la connecion a la base de datos*/
func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	return session
}

/*Acá apuntamos a la coleción de peliculas*/
var MovieCollections = getSession().DB("cursoUdemyGo").C("Movies")

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hola")
}

func MovieList(w http.ResponseWriter, r *http.Request) {
	var results []Movie

	// err := MovieCollections.Find(nil).All(&results)
	err := MovieCollections.Find(nil).Sort("-_id").All(&results)

	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(results)
}

func MovieShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]
	// fmt.Fprintf(w, "Has cargado la pelicula %s", movieID)

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	objectID := bson.ObjectIdHex(movieID)

	results := Movie{}

	err := MovieCollections.FindId(objectID).One(&results)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(results)

}

func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var movieData Movie

	err := decoder.Decode(&movieData)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	log.Println(movieData)

	err = MovieCollections.Insert(movieData)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(200)
	json.NewEncoder(w)

}

func MovieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	objectID := bson.ObjectIdHex(movieID)

	decoder := json.NewDecoder(r.Body)

	var movieData Movie

	err := decoder.Decode(&movieData)

	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}

	defer r.Body.Close()

	document := bson.M{"_id": objectID}

	change := bson.M{"$set": movieData}

	err = MovieCollections.Update(document, change)

	if err != nil {
		w.WriteHeader(404)
		panic(err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(movieData)

}

func MovieDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movieID)
	err := MovieCollections.RemoveId(oid)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	responseMessage := ResponseMessage{Status: "successful", Message: "eliminado exitosamente"}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(responseMessage)

}
