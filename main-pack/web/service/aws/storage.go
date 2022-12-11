package aws

import (
	"context"
	"ctf_cup/models"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
)

const (
	UsersBucket = "users"
	NotesBucket = "notes"
)

type Storage interface {
	GetUserInfo(name string, ctx context.Context) (models.UserInfo, error)
	PutUser(userInfo models.UserInfo, ctx context.Context) error
	IsUserExist(login string) (bool, error)
	PutNotes(dir string, note models.Note, ctx context.Context) error
	GetNotes(dir string, name string, ctx context.Context) (string, error)
}

type S3Storage struct {
	s3 BucketClient
}

func NewStorage(s3 BucketClient) Storage {
	return &S3Storage{s3: s3}
}

func (s *S3Storage) GetUserInfo(name string, ctx context.Context) (models.UserInfo, error) {
	get, err := s.s3.Get(UsersBucket, name, ctx)
	if err != nil {
		return models.UserInfo{}, err
	}

	var userInfo models.UserInfo
	err = json.Unmarshal(get, &userInfo)
	if err != nil {
		return models.UserInfo{}, err
	}

	return userInfo, nil
}

func (s *S3Storage) PutUser(userInfo models.UserInfo, ctx context.Context) error {
	userInfoStr, err := json.Marshal(userInfo)
	if err != nil {
		glog.Errorf("can't marshal userinfo: %s", err)
		return err
	}

	_, err = s.s3.Put(UsersBucket, userInfo.Credentials.Login, userInfoStr, ctx)
	if err != nil {
		glog.Errorf("can't write user info: %s", err)
		return err
	}

	return nil
}

func (s *S3Storage) IsUserExist(login string) (bool, error) {
	exist, err := s.s3.Exist(UsersBucket, login)
	if err != nil {
		glog.Errorf("can't handle message: %s", err)
		return false, err
	}

	return exist, nil
}

func (s *S3Storage) PutNotes(dir string, note models.Note, ctx context.Context) error {
	path := fmt.Sprintf("%s%s", dir, note.Name)
	_, err := s.s3.Put(NotesBucket, path, []byte(note.Text), ctx)
	if err != nil {
		glog.Errorf("can't write note: %s", err)
		return err
	}

	return nil
}

func (s *S3Storage) GetNotes(dir string, name string, ctx context.Context) (string, error) {
	path := fmt.Sprintf("%s%s", dir, name)
	data, err := s.s3.Get(NotesBucket, path, ctx)

	if err != nil {
		glog.Errorf("can't write note: %s", err)
		return "", err
	}

	return string(data), nil
}
