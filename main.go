package main

import (
	"fmt"
	"log"
	"mongotest/executor"
	"net/http"
	"os"
	// "sync"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pingcap/errors"
)

func main() {
  // test test test
	users := map[string]string{
		......
	}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.AuthBasicMiddleware{
		Realm: "auth",
		Authenticator: func(userId string, password string) bool {
			value, exists := users[userId]
			if exists && password == value {
				//fmt.Printf("User: %s\n", userId)
				log.Println("==============================")
				log.Printf("exec mongo by: %s", userId)
				return true
			}
			return false
		},
	})
	router, err := rest.MakeRouter(
		rest.Post("/mongo", PostMongo),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8088", api.MakeHandler()))
}

type Mongo struct {
	Db  string
	Sql string
}

//type Response struct {
//	Status string
//	Error  string
//}

//var store = map[string]*Response{}
//var lock = sync.RWMutex{}

func PostMongo(w rest.ResponseWriter, r *rest.Request) {
	mongo := Mongo{}
	err := r.DecodeJsonPayload(&mongo)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if mongo.Db == "" {
		rest.Error(w, "db required", 400)
		return
	}
	if mongo.Sql == "" {
		rest.Error(w, "sql required", 400)
		return
	}

	log.Printf("exec mongo sql: %s , %s\n", mongo.Db, mongo.Sql)

	//lock.Lock()
	//store[mongo.Db] = &mongo
	//lock.Unlock()
	//w.WriteJson(&mongo)

  // test test test
	var execution = "/data/program/mongodb/bin/mongo"
	var addr = "127.0.0.1:27017"
	var user = "admin"
	var password = "SOh3TbYhx8ypJPxmt"

	d, err := executor.NewExecutor(execution, addr, user, password, mongo.Db, mongo.Sql)
	if err != nil {
		fmt.Printf("Create Mongo error %v\n", errors.ErrorStack(err))
		//os.Exit(1)
	}

	var f = os.Stdout

	//if len(*output) > 0 {
	f, err = os.OpenFile("./cicd.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open file error %v\n", errors.ErrorStack(err))
		//os.Exit(1)
	}
	//}

	defer f.Close()

	//response := Response{"0", ""}

	if err = d.MongoExec(f); err != nil {
		fmt.Printf("Exec error %v\n", errors.ErrorStack(err))
		//response.Status = "1"
		//response.Error = errors.ErrorStack(err)
		rest.Error(w, "You have an error in your SQL syntax", 400)
		//return
		//os.Exit(1)
	}

	//w.WriteJson(&response)
	//rest.Error(w, "", 200)

}
