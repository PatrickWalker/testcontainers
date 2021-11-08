package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/db", DBServer)
	http.HandleFunc("/redis", RedisServer)
	http.HandleFunc("/dep", DepServer)

	fmt.Println("Server started and listening")
	http.ListenAndServe(":8080", nil)
}

//HelloServer is a helloworld function i stole from the internet. I feel remorse.
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I'm getting requests")
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

//DBServer is a postgres function I stole from the internet. I cannot stop
func DBServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DB endpoint hit")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_SCHEMA"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error connecting %v \n", err)
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging %v \n", err)

		panic(err)
	}
	fmt.Println("DB endpoint success")

	fmt.Fprintf(w, "DB Connection Works")

}

//RedisServer is a redis function I stole from the internet. I will not stop
func RedisServer(w http.ResponseWriter, r *http.Request) {
	conn, err := redis.Dial("tcp", os.Getenv("REDIS_URI"))
	if err != nil {
		fmt.Printf("Error connecting redis %v \n", err)
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always
	// properly closed before exiting the main() function.
	defer conn.Close()

	// Send our command across the connection. The first parameter to
	conn.Do("SET", "this", "that")
	fmt.Println("Redis endpoint hit")

	fmt.Fprintf(w, "We wrote something. This works")
}

//DepServer talks to a http dependency we have and it'll all be ok
func DepServer(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%v/postcodes", os.Getenv("SERVICE_URL")))
	if err != nil {
		fmt.Printf("Error connecting to dep: %v \n", err)

		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading resp : %v \n", err)

		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	fmt.Println("Dep endpoint hit")

	fmt.Fprintf(w, sb)
}
