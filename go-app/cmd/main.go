package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name        string `bson:"Name"`
	Description string `bson:"Description" default="Undefined`
	Country     string `bson:"Country" default="Undefined`
	Occupation  string `bson:"Occupation" default="Undefined`
	BirthYear   int    `bson:"BirthYear" default="Undefined`
	DeathYear   int    `bson:"DeathYear" default="Undefined`
}

func get_random_person() Person {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("ages").Collection("ages_collection")
	//This is quite slow. Let's try with a random index. 1124531 }
	filter := bson.D{{"Index", rand.Intn(10001)}}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result Person
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func main() {
	log.Println("Starting the server!")
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/random-person", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			tmpl := template.Must(template.ParseFiles("./templates/fragments/person.html"))
			data := get_random_person()
			tmpl.Execute(w, data)
		}
	})
	err := http.ListenAndServe(":3000", nil)
	if err == http.ErrServerClosed {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
