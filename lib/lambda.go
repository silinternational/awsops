package lib

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func LambdaInvoke(awsSess *session.Session, functionName, payload string) (*lambda.InvokeOutput, error) {
	svc := lambda.New(awsSess)

	encodedPayload := base64.StdEncoding.EncodeToString([]byte(payload))

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        []byte(encodedPayload),
	}

	return svc.Invoke(input)
}
