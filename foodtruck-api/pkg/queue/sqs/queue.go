package queue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/samber/lo"
)

type queuer struct {
	client   *sqs.SQS
	qURL     string
	qName    string
	region   string
	endpoint *string
}

// endpoint argument is optional
func New(qName string, region string, endpoint *string) *queuer {
	return &queuer{
		qName:    qName,
		region:   region,
		endpoint: endpoint,
	}
}

func (q *queuer) Open() (err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg := aws.Config{
		Region: lo.ToPtr(q.region),
	}
	// if endpoint is not empty, we will use localstack
	if q.endpoint != nil && *q.endpoint != "" {
		cfg.Endpoint = q.endpoint
	}

	q.client = sqs.New(sess, &cfg)

	var url string
	// create the queue if it does not exist
	if url, err = q.getQueueURL(); err == nil {
		q.qURL = url
		return nil
	}

	if url, err = q.createQueue(); err == nil {
		q.qURL = url
	}

	return err
}

func (q *queuer) createQueue() (url string, err error) {
	if q.client == nil {
		return url, errors.New("client is not initialized with Open().")
	}

	result, err := q.client.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &q.qName,
		Attributes: map[string]*string{
			"DelaySeconds":           lo.ToPtr("60"),
			"MessageRetentionPeriod": lo.ToPtr("86400"),
		},
	})

	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}

func (q *queuer) getQueueURL() (url string, err error) {
	if q.client == nil {
		return url, errors.New("client is not initialized with Open().")
	}

	result, err := q.client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &q.qName,
	})

	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}

// Enqueue a message.
// metadata argument is optional.
func (q *queuer) Enqueue(msg string, metadata map[string]string) error {
	if q.client == nil {
		return errors.New("client is not initialized with Open().")
	}

	var msgAttributes = make(map[string]*sqs.MessageAttributeValue)
	if metadata == nil {
		metadata = make(map[string]string)
	}

	for k, v := range metadata {
		msgAttributes[k] = &sqs.MessageAttributeValue{
			DataType:    lo.ToPtr("String"),
			StringValue: lo.ToPtr(v),
		}
	}

	byts, _ := json.Marshal(metadata)
	fmt.Println("message attributes", string(byts))

	_, err := q.client.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: msgAttributes,
		DelaySeconds:      lo.ToPtr(int64(10)),
		MessageBody:       lo.ToPtr(msg),
		QueueUrl:          &q.qURL,
	})
	return err
}

// Dequeue 1 message at a time.
//
// filterMetadata is optional.
// AWS SQS allows the receiver of the message to use the filter attributes to decide how to handle the message without having to first process the message body.
func (q *queuer) Dequeue(filterMetadata ...string) (metadata map[string]string, msgBody string, err error) {
	metadata = make(map[string]string)

	if q.client == nil {
		return metadata, msgBody, errors.New("client is not initialized with Open().")
	}

	var filters = []*string{lo.ToPtr(sqs.QueueAttributeNameAll)}
	if filterMetadata != nil {
		filters = lo.Map(filterMetadata, func(attr string, index int) *string {
			return &attr
		})
	}

	res, err := q.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			lo.ToPtr(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: filters,
		QueueUrl:              &q.qURL,
		MaxNumberOfMessages:   lo.ToPtr(int64(1)),
		VisibilityTimeout:     lo.ToPtr(int64(10)),
	})

	if len(res.Messages) > 0 {
		for k, v := range res.Messages[0].MessageAttributes {
			metadata[k] = *v.StringValue
		}
		msgBody = *res.Messages[0].Body
	}

	return metadata, msgBody, err
}
