package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

const (
	smsRegion = "us-east-1"
)

var (
	smsARN     = os.Getenv("GNOT_TOPIC_ARN")
	smsProfile = os.Getenv("GNOT_PROFILE")
)

func sendSMS(text, subject string) error {
	os.Setenv("AWS_PROFILE", smsProfile)
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := sns.New(sess, aws.NewConfig().WithRegion(smsRegion))

	in := &sns.PublishInput{
		Message: aws.String(text),
		// MessageAttributes: map[string]*sns.MessageAttributeValue{
		// 	"Key": {
		// 		DataType:    aws.String("String"),
		// 		BinaryValue: []byte("PAYLOAD"),
		// 		StringValue: aws.String("String"),
		// 	},
		// },
		Subject:  aws.String(subject),
		TopicArn: aws.String(smsARN),
	}

	_, err = svc.Publish(in)
	return err
}
