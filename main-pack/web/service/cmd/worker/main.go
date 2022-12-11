package main

import (
	"context"
	"ctf_cup/aws"
	"fmt"
	"github.com/golang/glog"
	"os"
	"time"
)

const host = "http://localstack:4566"

func main() {
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

	flag := os.Getenv("CTFCUP_FLAG")
	for {
		_, err = s3.Put(aws.NotesBucket, "secretNote", []byte(flag), context.Background())
		if err != nil {
			glog.Error(err.Error())
			<-time.Tick(time.Second * 5)
			continue
		}

		break
	}

	storage := aws.NewStorage(s3)
	StartWorker(queue, storage)
}
