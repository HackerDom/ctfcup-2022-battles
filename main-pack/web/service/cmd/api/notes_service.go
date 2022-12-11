package main

import (
	"context"
	"ctf_cup/aws"
	"ctf_cup/models"
	"errors"
	"github.com/golang/glog"
	"time"
)

type NotesService interface {
	CreateUser(credentials models.Credentials, ctx context.Context) error
	WriteNotes(note models.Note, credentials models.Credentials, ctx context.Context) error
	GetNotes(path string, credentials models.Credentials, ctx context.Context) (string, error)
}

type NotesServiceImpl struct {
	queue   aws.Queue
	users   map[string]models.Credentials
	storage aws.Storage
}

func NewNotesService(queue aws.Queue, storage aws.Storage) NotesService {
	return &NotesServiceImpl{
		queue:   queue,
		users:   map[string]models.Credentials{},
		storage: storage,
	}
}

func (service *NotesServiceImpl) CreateUser(credentials models.Credentials, ctx context.Context) error {
	exist, err := service.storage.IsUserExist(credentials.Login)
	if err != nil {
		glog.Error(err)
		return err
	}

	if exist {
		return errors.New("user exist")
	}

	msg := models.CreateUser{Credentials: credentials}.ToMessage()

	msgId, err := service.queue.Push(msg.ToString(), ctx)
	if err != nil {
		glog.Errorf("Send failed: %s", err)
		return err
	}

	for range time.Tick(5 * time.Second) {
		inQueue, err := service.queue.InQueue(msgId)
		if err != nil {
			return err
		}

		if !inQueue {
			service.users[credentials.Login] = credentials
			break
		}
	}

	glog.Infof("Successfully add: %+v", credentials)

	return nil
}

func (service *NotesServiceImpl) WriteNotes(note models.Note, credentials models.Credentials, ctx context.Context) error {
	err := service.unsureAuthenticated(credentials)
	if err != nil {
		return err
	}

	message := models.WriteNote{
		Note: note,
		Dir:  "",
	}

	info, err := service.storage.GetUserInfo(credentials.Login, ctx)
	message.Dir = info.HomeDir
	send, err := service.queue.Push(message.ToMessage().ToString(), ctx)
	if err != nil {
		glog.Errorf("Send failed: %s", err)
		return err
	}

	glog.Infof("Successfully sent: %s", send)
	return nil
}

func (service *NotesServiceImpl) unsureAuthenticated(credentials models.Credentials) error {
	expected, ok := service.users[credentials.Login]

	if !ok {
		return errors.New("user not exist")
	}

	if credentials.Password != expected.Password {
		return errors.New("incorrect password")
	}
	return nil
}

func (service *NotesServiceImpl) GetNotes(name string, credentials models.Credentials, ctx context.Context) (string, error) {
	err := service.unsureAuthenticated(credentials)
	if err != nil {
		return "", err
	}

	info, err := service.storage.GetUserInfo(credentials.Login, ctx)
	note, err := service.storage.GetNotes(info.HomeDir, name, ctx)
	if err != nil {
		return "", err
	}

	return note, nil
}
