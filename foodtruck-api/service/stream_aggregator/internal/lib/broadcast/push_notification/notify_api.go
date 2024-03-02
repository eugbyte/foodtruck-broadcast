package notify

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/samber/lo"
)

type notifyLambda struct {
	client     *lambda.Lambda
	lambdaName string
}

func NewNotificationAPI(region string, lambdaName string) *notifyLambda {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &notifyLambda{
		client:     lambda.New(sess, &aws.Config{Region: lo.ToPtr(region)}),
		lambdaName: lambdaName,
	}
}

func (l *notifyLambda) Post(payload map[string]any) (respBody []byte, err error) {
	byts, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, err
	}

	input := lambda.InvokeInput{FunctionName: &l.lambdaName, Payload: byts}
	resp, err := l.client.Invoke(&input)
	return resp.Payload, err
}
