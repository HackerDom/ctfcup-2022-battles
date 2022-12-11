package main

import (
	"ctf_cup/aws"
	"ctf_cup/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
)

const host = "http://localstack:4566"

var service NotesService

func main() {
	flag.Parse()

	config := aws.Config{
		Address: host,
		Region:  "eu-west-1",
		Profile: "localstack",
		ID:      "empty",
		Secret:  "empty",
	}
	queueUrl := fmt.Sprintf("%s/%s", host, aws.CommandQueue)
	queue, err := aws.NewSQSQueue(config, queueUrl)
	if err != nil {
		glog.Fatal(err)
	}

	s3 := aws.NewS3(config)
	storage := aws.NewStorage(s3)

	service = NewNotesService(queue, storage)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", HandleRegisterRequest)
	mux.HandleFunc("/putNote", HandlePutNoteRequest)
	mux.HandleFunc("/getNote", HandleGetNoteRequest)

	glog.Infof("Starting notes service")
	http.ListenAndServe(":8080", mux)
}

func HandleGetNoteRequest(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var putNoteRequest models.GetNoteRequest
	err = json.Unmarshal(data, &putNoteRequest)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	note, err := service.GetNotes(putNoteRequest.Name, putNoteRequest.Credentials, request.Context())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.WriteHeader(200)
	writer.Write([]byte(note))
}

func HandlePutNoteRequest(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var putNoteRequest models.PutNoteRequest
	err = json.Unmarshal(data, &putNoteRequest)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	err = service.WriteNotes(putNoteRequest.Note, putNoteRequest.Credentials, request.Context())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.WriteHeader(200)
}

func HandleRegisterRequest(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	var credentials models.Credentials
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	err = service.CreateUser(credentials, request.Context())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.WriteHeader(200)
}
