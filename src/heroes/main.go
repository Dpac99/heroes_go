package main

import(
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"
)

var collection *mongo.Collection

type Hero struct{
	Id int 
	Name string 
}


func Home(w http.ResponseWriter, r* http.Request){
	fmt.Fprintf(w, "Welcome to heroes!")
}

func GetHeroes(w http.ResponseWriter, r* http.Request){
	var result []Hero
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item Hero
		err := cur.Decode(&item)
		if err != nil { log.Fatal(err) }
		result = append(result, item)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}

func PostHero(w http.ResponseWriter, r* http.Request){
	decoder := json.NewDecoder(r.Body)
	var h Hero
	err:=decoder.Decode(&h)
	if err != nil{
		fmt.Fprintf(w, "Error with payload")
		return
	}
	_, err = collection.InsertOne(context.TODO(), h)
	if err != nil {
    	log.Fatal(err)
	}
	json.NewEncoder(w).Encode(h)
}

func GetHero(w http.ResponseWriter, r* http.Request){
	id, _ := strconv.Atoi(mux.Vars(r)["id"]);
	filter := bson.D{{"id", id}}
	var h Hero
	err := collection.FindOne(context.TODO(), filter).Decode(&h)
	if err != nil {
    	fmt.Println("No Heroes found with such ID")
	}
	json.NewEncoder(w).Encode(h)
}

func DeleteHero(w http.ResponseWriter, r* http.Request){
	id, _ := strconv.Atoi(mux.Vars(r)["id"]);
	filter := bson.D{{"id", id}}
	var h Hero
	err := collection.FindOne(context.TODO(), filter).Decode(&h)
	if err != nil {
    	fmt.Println("No Heroes found with such ID")
	}
	_,err = collection.DeleteOne(context.TODO(), filter)
	json.NewEncoder(w).Encode(h)
}

func UpdateHero(w http.ResponseWriter, r* http.Request){
	decoder := json.NewDecoder(r.Body)
	id, _ := strconv.Atoi(mux.Vars(r)["id"]);
	var hero Hero
	err := decoder.Decode(&hero)

	if err != nil{
		fmt.Fprintf(w, "Error with payload")
		return
	}

	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"name", hero.Name},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
    	fmt.Fprintf(w, "No heroes found")
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&hero)

	json.NewEncoder(w).Encode(hero)

}

func main(){

	ctx1, c1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer c1()
	client, err := mongo.Connect(ctx1, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{ log.Fatal(err)}
	ctx2, c2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer c2()
	err = client.Ping(ctx2, readpref.Primary())
	collection = client.Database("local").Collection("heroes")	

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/hero", GetHeroes).Methods("GET")
	r.HandleFunc("/hero", PostHero).Methods("POST")
	r.HandleFunc("/hero/{id}", GetHero).Methods("GET")
	r.HandleFunc("/hero/{id}", DeleteHero).Methods("DELETE")
	r.HandleFunc("/hero/{id}", UpdateHero).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}