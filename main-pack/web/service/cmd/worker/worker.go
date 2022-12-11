package main

import (
	"context"
	"ctf_cup/aws"
	"ctf_cup/models"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"time"
)

const batchSize = 50

type daemonImpl struct {
	queue   aws.Queue
	storage aws.Storage
}

func StartWorker(queue aws.Queue, storage aws.Storage) {
	worker := daemonImpl{
		queue:   queue,
		storage: storage,
	}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := worker.pull(ctx)
		if err != nil {
			glog.Errorf("task pull failed. Reason: %s", err)
			time.Sleep(1 * time.Second)
		}
		cancel()
	}
}

func (worker *daemonImpl) pull(ctx context.Context) error {
	glog.Info("Start pulling messages")

	var batch []aws.Message
	for i := 0; i < batchSize; i++ {
		msgs, err := worker.queue.Pull(ctx)
		if err != nil {
			glog.Errorf("can't pull: %s", err)
			return err
		}

		batch = append(batch, msgs...)

		if len(msgs) == 0 {
			break
		}
	}

	for _, msg := range batch {
		body := msg.Body
		worker.handleMessage(ctx, body)
		_ = worker.queue.Delete(msg)
	}

	return nil
}

func (worker *daemonImpl) handleMessage(ctx context.Context, body string) {
	var message models.Message
	err := json.Unmarshal([]byte(body), &message)
	if err != nil {
		glog.Errorf("can't handle message: %s", err)
		return
	}

	switch message.Type {
	case models.CreateUserMessage:
		worker.HandleCreateUser(message, ctx)
	case models.WriteNoteMessage:
		worker.HandleWriteNote(message, ctx)
	default:
		glog.Errorf("unknown message type: %+v", message)
	}
}
func (worker *daemonImpl) HandleWriteNote(message models.Message, ctx context.Context) {
	var writeNoteMsg models.WriteNote
	err := json.Unmarshal(message.Data, &writeNoteMsg)
	if err != nil {
		glog.Errorf("can't handle message: %s", err)
		return
	}

	err = worker.storage.PutNotes(writeNoteMsg.Dir, writeNoteMsg.Note, ctx)
	if err != nil {
		glog.Errorf("can't write note: %s", err)
	}
}

func (worker *daemonImpl) HandleCreateUser(message models.Message, ctx context.Context) {
	var createUserMsg models.CreateUser
	err := json.Unmarshal(message.Data, &createUserMsg)
	if err != nil {
		glog.Errorf("can't handle message: %s", err)
		return
	}

	exist, err := worker.storage.IsUserExist(createUserMsg.Credentials.Login)
	if err != nil {
		glog.Errorf("can't handle message: %s", err)
		return
	}

	if exist {
		glog.Errorf("user exist: %s", err)
		return
	}

	userInfo := models.UserInfo{
		Credentials: createUserMsg.Credentials,
		HomeDir:     uuid.New().String(),
	}

	err = worker.storage.PutUser(userInfo, ctx)
	if err != nil {
		glog.Errorf("can't write user info: %s", err)
	}
}
