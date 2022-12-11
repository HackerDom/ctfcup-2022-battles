package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

const connectionString = "postgresql://postgres/secrets?user=postgres&password=empty&sslmode=disable"

var db *sqlx.DB
var auth map[string]*User

type User struct {
	Name  string `db:"name"`
	Token string `db:"token"`
}

func main() {
	flag.Parse()
	var err error
	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		glog.Fatal(err.Error())
		return
	}

	Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/createUser", createUser)
	mux.HandleFunc("/getSecrets", getSecrets)
	mux.HandleFunc("/addSecret", addSecret)

	glog.Infof("Starting secret service")
	http.ListenAndServe(":8080", mux)
}

func Init() {
	auth = map[string]*User{}
	users := GetUsers()
	for _, user := range users {
		auth[user.Name] = &user
	}
}

func addSecret(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	token := request.URL.Query().Get("token")
	secret := request.URL.Query().Get("secret")

	user, ok := auth[name]
	if !ok {
		writer.WriteHeader(404)
		writer.Write([]byte("user not found"))
		return
	}

	if user.Token != token {
		writer.WriteHeader(403)
		writer.Write([]byte("user not found"))
		return
	}

	err := AddSecret(name, secret)
	if err != nil {
		glog.Error(err)
	}
}

func getSecrets(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	token := request.URL.Query().Get("token")

	user, ok := auth[name]
	if !ok {
		writer.WriteHeader(404)
		writer.Write([]byte("user not found"))
		return
	}

	if user.Token != token {
		writer.WriteHeader(403)
		writer.Write([]byte("user not found"))
		return
	}

	userSecrets := GetSecrets(name)
	writer.WriteHeader(200)
	for _, secret := range userSecrets {
		writer.Write([]byte(secret + "\n"))
	}
}

func createUser(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	token := uuid.New().String()

	err := AddUser(name, token)
	if err != nil {
		writer.WriteHeader(400)
		glog.Error(err)
		return
	}

	auth[name] = &User{
		Name:  name,
		Token: token,
	}

	writer.WriteHeader(200)
	writer.Write([]byte(token))
}
