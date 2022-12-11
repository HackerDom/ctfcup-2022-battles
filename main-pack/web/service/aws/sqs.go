package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/glog"
	"time"
)

const messageStateBucket = "messages"

const (
	CommandQueue = "000000000000/commandQueue"
)

func GetCommandQueueUrl(host string) string {
	return fmt.Sprintf("%s/%s", host, CommandQueue)
}

type Message struct {
	Body    string
	receipt string
}

type Queue interface {
	Push(message string, ctx context.Context) (string, error)
	Pull(ctx context.Context) ([]Message, error)
	Delete(msg Message) error
	InQueue(messageId string) (bool, error)
}

type SqsClient struct {
	client   *sqs.SQS
	queueUrl string
	s3       BucketClient
}

func NewSQSQueue(config Config, queueURL string) (Queue, error) {
	s, err := newAWSSession(config)
	if err != nil {
		return nil, err
	}

	return &SqsClient{
		client:   sqs.New(s),
		queueUrl: queueURL,
		s3:       NewS3(config),
	}, nil
}

func (s *SqsClient) Push(message string, ctx context.Context) (string, error) {

	res, err := s.client.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(s.queueUrl),
	})
	if err != nil {
		glog.Errorf("send: %s", err)
		return "", err
	}

	_, err = s.s3.Put(messageStateBucket, *res.MessageId, []byte("InQueue"), ctx)
	if err != nil {
		glog.Errorf("can't change message state: %s", err)
	}

	return *res.MessageId, nil
}

func (s *SqsClient) InQueue(messageId string) (bool, error) {
	return s.s3.Exist(messageStateBucket, messageId)
}

func (s *SqsClient) Pull(ctx context.Context) ([]Message, error) {
	res, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueUrl),
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	})
	if err != nil {
		glog.Errorf("receive: %s", err)
		return nil, err
	}

	if len(res.Messages) == 0 {
		return nil, nil
	}

	var result []Message
	for _, message := range res.Messages {
		if message.Body == nil {
			continue
		}

		err = s.s3.Delete(messageStateBucket, *message.MessageId, ctx)
		if err != nil {
			glog.Errorf("can't change message state: %s", err)
		}

		result = append(result, Message{
			Body:    *message.Body,
			receipt: *message.ReceiptHandle,
		})
		err := s.hide(*message.ReceiptHandle, time.Second*30)
		if err != nil {
			glog.Errorf("can't hide: %s", err)
			return nil, err
		}
	}

	return result, nil
}

func (s *SqsClient) Delete(msg Message) error {
	glog.Infof("deleting message: %s", msg.receipt)

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueUrl),
		ReceiptHandle: aws.String(msg.receipt),
	}
	_, err := s.client.DeleteMessage(params)
	return err
}

func (s *SqsClient) hide(receipt string, interval time.Duration) (err error) {
	glog.Info("hiding message")

	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(s.queueUrl),
		ReceiptHandle:     aws.String(receipt),
		VisibilityTimeout: aws.Int64(int64(interval.Seconds())),
	}

	_, err = s.client.ChangeMessageVisibility(params)
	return
}
